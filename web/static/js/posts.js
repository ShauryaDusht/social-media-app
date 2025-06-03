// Base API URL
const API_URL = '/api';

// Get token from local storage
const token = localStorage.getItem('token');
const user = JSON.parse(localStorage.getItem('user'));

// Initialize posts page
async function initPostsPage() {
    // Set up post form submission
    const postForm = document.getElementById('post-form');
    postForm.addEventListener('submit', createPost);
    
    // Load posts
    await loadPosts();
}

// Create a new post
async function createPost(e) {
    e.preventDefault();
    const content = document.getElementById('post-content').value;
    const imageUrl = document.getElementById('image-url').value;
    
    try {
        const response = await fetch(`${API_URL}/posts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                content,
                image_url: imageUrl
            })
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to create post');
        }
        
        // Clear form
        document.getElementById('post-content').value = '';
        document.getElementById('image-url').value = '';
        
        // Reload posts
        await loadPosts();
    } catch (error) {
        console.error('Error creating post:', error);
        alert(error.message);
    }
}

// Load posts
async function loadPosts() {
    const postsContainer = document.getElementById('posts-list');
    postsContainer.innerHTML = '<div class="loading">Loading posts...</div>';
    
    try {
        const response = await fetch(`${API_URL}/posts`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to load posts');
        }
        
        if (data.data.length === 0) {
            postsContainer.innerHTML = '<p>No posts yet. Be the first to post!</p>';
            return;
        }
        
        // Render posts
        postsContainer.innerHTML = '';
        data.data.forEach(post => {
            postsContainer.appendChild(createPostElement(post));
        });
    } catch (error) {
        console.error('Error loading posts:', error);
        postsContainer.innerHTML = `<p>Error loading posts: ${error.message}</p>`;
    }
}

// Create post element
function createPostElement(post) {
    const postElement = document.createElement('div');
    postElement.className = 'post-card';
    postElement.dataset.id = post.id;
    
    const isLiked = post.is_liked;
    const isOwnPost = post.user.id === user.id;
    
    postElement.innerHTML = `
        <div class="post-header">
            <img src="${post.user.avatar || '/static/img/default-avatar.png'}" alt="${post.user.username}" class="post-avatar">
            <span class="post-user">${post.user.first_name} ${post.user.last_name}</span>
            <span class="post-time">${new Date(post.created_at).toLocaleString()}</span>
        </div>
        <div class="post-content">${post.content}</div>
        ${post.image_url ? `<img src="${post.image_url}" alt="Post image" class="post-image">` : ''}
        <div class="post-actions">
            <div class="post-action ${isLiked ? 'liked' : ''}" onclick="toggleLike(${post.id}, ${isLiked})">
                ❤ <span class="like-count">${post.like_count}</span> Likes
            </div>
            ${isOwnPost ? `
                <div class="post-action" onclick="editPost(${post.id})">✏️ Edit</div>
                <div class="post-action" onclick="deletePost(${post.id})">🗑️ Delete</div>
            ` : ''}
        </div>
    `;
    
    return postElement;
}

// Toggle like on a post
async function toggleLike(postId, isLiked) {
    try {
        const method = isLiked ? 'DELETE' : 'POST';
        const url = isLiked ? `${API_URL}/likes/${postId}` : `${API_URL}/likes`;
        
        const response = await fetch(url, {
            method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: isLiked ? null : JSON.stringify({ post_id: postId })
        });
        
        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || 'Failed to toggle like');
        }
        
        // Reload posts to update like count
        await loadPosts();
    } catch (error) {
        console.error('Error toggling like:', error);
        alert(error.message);
    }
}

// Delete a post
async function deletePost(postId) {
    if (!confirm('Are you sure you want to delete this post?')) {
        return;
    }
    
    try {
        const response = await fetch(`${API_URL}/posts/${postId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || 'Failed to delete post');
        }
        
        // Reload posts
        await loadPosts();
    } catch (error) {
        console.error('Error deleting post:', error);
        alert(error.message);
    }
}

// Edit a post (simplified - in a real app, you'd use a modal or inline editing)
async function editPost(postId) {
    const newContent = prompt('Edit your post:');
    if (newContent === null) return; // User cancelled
    
    try {
        const response = await fetch(`${API_URL}/posts/${postId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                content: newContent
            })
        });
        
        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || 'Failed to update post');
        }
        
        // Reload posts
        await loadPosts();
    } catch (error) {
        console.error('Error updating post:', error);
        alert(error.message);
    }
}

// Initialize on page load
initPostsPage();