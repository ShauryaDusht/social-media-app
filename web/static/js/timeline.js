const API_URL = '/api';
const token = localStorage.getItem('token');
const user = JSON.parse(localStorage.getItem('user'));

async function initTimelinePage() {
    if (!token) {
        showLoginPrompt();
        return;
    }
    
    await loadTimelinePosts();
}

function showLoginPrompt() {
    const postsContainer = document.getElementById('timeline-posts');
    const titleElement = document.getElementById('timeline-title');
    
    if (titleElement) {
        titleElement.textContent = 'Welcome to Social Media App';
    }
    
    if (postsContainer) {
        postsContainer.innerHTML = `
            <div style="text-align: center; padding: 40px; background: #f8f9fa; border-radius: 8px; margin: 20px 0;">
                <h3>Join the Community!</h3>
                <p>Please login to view your personalized timeline and interact with posts.</p>
                <a href="/login" style="display: inline-block; background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 10px;">Login</a>
                <a href="/signup" style="display: inline-block; background: #28a745; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 10px;">Sign Up</a>
            </div>
        `;
    }
}

async function loadTimelinePosts() {
    const postsContainer = document.getElementById('timeline-posts');
    
    if (!postsContainer) {
        return;
    }
    
    postsContainer.innerHTML = '<div class="loading">Loading your timeline...</div>';
    
    try {
        const response = await fetch(`${API_URL}/posts`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            showLoginPrompt();
            return;
        }
        
        if (!response.ok) {
            throw new Error('Failed to load timeline posts');
        }
        
        const data = await response.json();
        const posts = data.data || data || [];
        
        if (posts.length === 0) {
            postsContainer.innerHTML = `
                <div style="text-align: center; padding: 40px; background: #f8f9fa; border-radius: 8px; margin: 20px 0;">
                    <h3>Your Timeline is Empty</h3>
                    <p>Follow some users or create your first post to see content here.</p>
                    <a href="/posts" style="display: inline-block; background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 10px;">Create Post</a>
                    <a href="/users" style="display: inline-block; background: #6c757d; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 10px;">Find Users</a>
                </div>
            `;
            return;
        }
        
        postsContainer.innerHTML = '';
        posts.forEach(post => {
            const postElement = createTimelinePostElement(post);
            postsContainer.appendChild(postElement);
        });
        
    } catch (error) {
        postsContainer.innerHTML = `
            <div style="text-align: center; padding: 40px; background: #ffe6e6; border-radius: 8px; margin: 20px 0; color: #d63384;">
                <h3>Error Loading Timeline</h3>
                <p>${error.message}</p>
                <button onclick="loadTimelinePosts()" style="background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer;">Try Again</button>
            </div>
        `;
    }
}

function createTimelinePostElement(post) {
    const postElement = document.createElement('div');
    postElement.className = 'post-card';
    postElement.dataset.id = post.id;
    
    const isLiked = post.is_liked;
    const isOwnPost = user && post.user && post.user.id === user.id;
    
    let dateString = 'Unknown date';
    try {
        if (post.created_at) {
            dateString = new Date(post.created_at).toLocaleString();
        }
    } catch (e) {
        dateString = 'Unknown date';
    }
    
    postElement.innerHTML = `
        <div class="post-header">
            <img src="${post.user?.avatar || '/static/img/default-avatar.png'}" alt="${post.user?.username || 'User'}" class="post-avatar">
            <div class="post-user-info">
                <span class="post-user">${post.user?.first_name || 'Anonymous'} ${post.user?.last_name || ''}</span>
                <span class="post-username">@${post.user?.username || 'unknown'}</span>
            </div>
            <span class="post-time">${dateString}</span>
        </div>
        <div class="post-content">${post.content || 'No content'}</div>
        ${post.image_url ? `<img src="${post.image_url}" alt="Post image" class="post-image">` : ''}
        <div class="post-actions">
            <div class="post-action ${isLiked ? 'liked' : ''}" onclick="toggleLike(${post.id}, ${isLiked})">
                ❤ <span class="like-count">${post.like_count || 0}</span> Likes
            </div>
            <div class="post-action" onclick="viewProfile(${post.user?.id})">
                👤 View Profile
            </div>
            ${isOwnPost ? `
                <div class="post-action" onclick="editPost(${post.id})">✏️ Edit</div>
                <div class="post-action" onclick="deletePost(${post.id})">🗑️ Delete</div>
            ` : ''}
        </div>
    `;
    
    return postElement;
}

async function toggleLike(postId, isLiked) {
    if (!token) {
        alert('Please login to like posts');
        return;
    }
    
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
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to toggle like');
        }
        
        // Update UI without reloading all posts
        const postElement = document.querySelector(`.post-card[data-id="${postId}"]`);
        if (postElement) {
            const likeAction = postElement.querySelector('.post-action');
            const likeCount = postElement.querySelector('.like-count');
            
            if (likeAction && likeCount) {
                if (isLiked) {
                    likeAction.classList.remove('liked');
                    likeCount.textContent = Math.max(0, parseInt(likeCount.textContent) - 1);
                } else {
                    likeAction.classList.add('liked');
                    likeCount.textContent = parseInt(likeCount.textContent) + 1;
                }
            }
        }
    } catch (error) {
        console.error('Error toggling like:', error);
        alert(error.message);
    }
}

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
            throw new Error('Failed to delete post');
        }
        
        await loadTimelinePosts();
    } catch (error) {
        alert(error.message);
    }
}

async function editPost(postId) {
    const newContent = prompt('Edit your post:');
    if (newContent === null) {
        return;
    }
    
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
            throw new Error('Failed to update post');
        }
        
        await loadTimelinePosts();
    } catch (error) {
        alert(error.message);
    }
}

function viewProfile(userId) {
    if (userId) {
        window.location.href = `/profile?id=${userId}`;
    }
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initTimelinePage);
} else {
    initTimelinePage();
}