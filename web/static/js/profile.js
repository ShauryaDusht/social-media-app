// Base API URL
const API_URL = '/api';

// Get token from local storage
const token = localStorage.getItem('token');
const currentUser = JSON.parse(localStorage.getItem('user'));
const currentUserId = currentUser ? currentUser.id : null;

// Get profile user ID from URL or use current user ID
let profileUserId = currentUserId;
const urlParams = new URLSearchParams(window.location.search);
if (urlParams.has('id')) {
    profileUserId = parseInt(urlParams.get('id'));
}

// Initialize profile page
async function initProfilePage() {
    // Set up tab switching
    const tabButtons = document.querySelectorAll('.tab-btn');
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            // Remove active class from all buttons and hide all content
            tabButtons.forEach(btn => btn.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(content => {
                content.classList.add('hidden');
            });
            
            // Add active class to clicked button and show corresponding content
            button.classList.add('active');
            const tabId = button.dataset.tab;
            document.getElementById(`${tabId}-tab`).classList.remove('hidden');
        });
    });
    
    // Set up profile edit form
    const editForm = document.getElementById('edit-profile-form');
    editForm.addEventListener('submit', updateProfile);
    
    // Hide edit tab if viewing another user's profile
    if (profileUserId !== currentUserId) {
        document.querySelector('[data-tab="edit"]').style.display = 'none';
        document.getElementById('follow-container').style.display = 'block';
    }
    
    // Load profile data
    await loadProfile();
    
    // Load user posts
    await loadUserPosts();
}

// Load profile data
async function loadProfile() {
    try {
        let url = `${API_URL}/users/profile`;
        if (profileUserId !== currentUserId) {
            url = `${API_URL}/users/${profileUserId}`;
        }
        
        const response = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        // Check if response is empty
        const text = await response.text();
        const data = text ? JSON.parse(text) : {};
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to load profile');
        }
        
        const profile = data.data;
        
        // Update profile display
        document.getElementById('user-name').textContent = `${profile.first_name} ${profile.last_name}`;
        document.getElementById('user-username').textContent = `@${profile.username}`;
        document.getElementById('user-bio').textContent = profile.bio || 'No bio yet';
        
        if (profile.avatar) {
            document.getElementById('avatar').src = profile.avatar;
        }
        
        // Populate edit form if it's the current user's profile
        if (profileUserId === currentUserId) {
            document.getElementById('first_name').value = profile.first_name;
            document.getElementById('last_name').value = profile.last_name;
            document.getElementById('bio').value = profile.bio || '';
            document.getElementById('avatar-url').value = profile.avatar || '';
        }
        
        // Load followers and following counts
        await loadFollowStats(profile.id);
    } catch (error) {
        console.error('Error loading profile:', error);
        alert(`Error loading profile: ${error.message}`);
    }
}

// Load follow statistics
async function loadFollowStats(userId) {
    try {
        // Get followers
        const followersResponse = await fetch(`${API_URL}/follows/followers/${userId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const followersData = await followersResponse.json();
        
        if (!followersResponse.ok) {
            throw new Error(followersData.error || 'Failed to load followers');
        }
        
        // Get following
        const followingResponse = await fetch(`${API_URL}/follows/following/${userId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const followingData = await followingResponse.json();
        
        if (!followingResponse.ok) {
            throw new Error(followingData.error || 'Failed to load following');
        }
        
        // Update stats display
        document.getElementById('followers-count').textContent = followersData.data.length;
        document.getElementById('following-count').textContent = followingData.data.length;
        
        // Check if current user is following this profile
        if (profileUserId !== currentUserId) {
            const isFollowing = followersData.data.some(follow => follow.user.id === currentUserId);
            // Use the updateFollowButton function from follows.js
            if (typeof updateFollowButton === 'function') {
                updateFollowButton(isFollowing);
            }
        }
    } catch (error) {
        console.error('Error loading follow stats:', error);
    }
}

// These functions have been moved to follows.js

// Load user posts
async function loadUserPosts() {
    const postsContainer = document.getElementById('user-posts');
    postsContainer.innerHTML = '<div class="loading">Loading posts...</div>';
    
    try {
        const response = await fetch(`${API_URL}/posts/user/${user.id}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to load posts');
        }
        
        // Update posts count
        document.getElementById('posts-count').textContent = data.data.length;
        
        if (data.data.length === 0) {
            postsContainer.innerHTML = '<p>No posts yet.</p>';
            return;
        }
        
        // Render posts
        postsContainer.innerHTML = '';
        data.data.forEach(post => {
            const postElement = document.createElement('div');
            postElement.className = 'post-card';
            
            postElement.innerHTML = `
                <div class="post-header">
                    <span class="post-time">${new Date(post.created_at).toLocaleString()}</span>
                </div>
                <div class="post-content">${post.content}</div>
                ${post.image_url ? `<img src="${post.image_url}" alt="Post image" class="post-image">` : ''}
                <div class="post-actions">
                    <div class="post-action">
                        ❤ ${post.like_count} Likes
                    </div>
                    <div class="post-action" onclick="editPost(${post.id})">✏️ Edit</div>
                    <div class="post-action" onclick="deletePost(${post.id})">🗑️ Delete</div>
                </div>
            `;
            
            postsContainer.appendChild(postElement);
        });
    } catch (error) {
        console.error('Error loading user posts:', error);
        postsContainer.innerHTML = `<p>Error loading posts: ${error.message}</p>`;
    }
}

// Update profile
async function updateProfile(e) {
    e.preventDefault();
    
    const firstName = document.getElementById('first_name').value;
    const lastName = document.getElementById('last_name').value;
    const bio = document.getElementById('bio').value;
    const avatar = document.getElementById('avatar-url').value;
    const password = document.getElementById('new-password').value;
    
    const updateData = {
        first_name: firstName,
        last_name: lastName,
        bio,
        avatar
    };
    
    // Only include password if it was provided
    if (password) {
        updateData.password = password;
    }
    
    try {
        const response = await fetch(`${API_URL}/users/profile`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(updateData)
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to update profile');
        }
        
        alert('Profile updated successfully!');
        
        // Update local storage user data
        const updatedUser = {
            ...user,
            first_name: firstName,
            last_name: lastName,
            bio,
            avatar
        };
        localStorage.setItem('user', JSON.stringify(updatedUser));
        
        // Reload profile
        await loadProfile();
        
        // Switch back to posts tab
        document.querySelector('[data-tab="posts"]').click();
    } catch (error) {
        console.error('Error updating profile:', error);
        alert(`Error updating profile: ${error.message}`);
    }
}

// Edit a post
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
        
        // Reload user posts
        await loadUserPosts();
    } catch (error) {
        console.error('Error updating post:', error);
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
        
        // Reload user posts
        await loadUserPosts();
    } catch (error) {
        console.error('Error deleting post:', error);
        alert(error.message);
    }
}

// Initialize on page load
initProfilePage();