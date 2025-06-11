const API_URL = '/api';

const token = localStorage.getItem('token');
const currentUser = JSON.parse(localStorage.getItem('user'));
const currentUserId = currentUser ? currentUser.id : null;

let profileUserId = currentUserId;
const urlParams = new URLSearchParams(window.location.search);
if (urlParams.has('id')) {
    profileUserId = parseInt(urlParams.get('id'));
}


async function initProfilePage() {
    const tabButtons = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');
    
    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            tabButtons.forEach(btn => btn.classList.remove('active'));
            tabContents.forEach(content => {
                content.style.display = 'none';
            });
            
            button.classList.add('active');
            const tabId = button.dataset.tab;
            const tabElement = document.getElementById(`${tabId}-tab`);
            
            if (tabElement) {
                tabElement.style.display = 'block';
            }
        });
    });
    
    document.getElementById('posts-tab').style.display = 'block';
    document.getElementById('edit-profile-tab').style.display = 'none';
    
    const editForm = document.getElementById('edit-profile-form');
    if (editForm) {
        editForm.addEventListener('submit', updateProfile);
    }
    
    if (profileUserId !== currentUserId) {
        const editTabBtn = document.querySelector('[data-tab="edit-profile"]');
        if (editTabBtn) {
            editTabBtn.style.display = 'none';
        }
        
        const followContainer = document.getElementById('follow-container');
        if (followContainer) {
            followContainer.style.display = 'block';
        }
    }
    
    await loadProfile();
    await loadUserPosts();
}

async function loadProfile() {
    try {
        let url = `/api/users/profile`;
        if (profileUserId !== currentUserId) {
            url = `/api/users/${profileUserId}`;
        }
        
        const response = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load profile');
        }
        
        const text = await response.text();
        const data = text ? JSON.parse(text) : {};
        const profile = data.data;
        
        const profileNameEl = document.getElementById('profile-name');
        const profileUsernameEl = document.getElementById('profile-username');
        const profileBioEl = document.getElementById('profile-bio');
        const profileAvatarEl = document.getElementById('profile-avatar');
        
        if (profileNameEl) {
            profileNameEl.textContent = `${profile.first_name || ''} ${profile.last_name || ''}`.trim();
        }
        if (profileUsernameEl) {
            profileUsernameEl.textContent = `@${profile.username || 'unknown'}`;
        }
        if (profileBioEl) {
            profileBioEl.textContent = profile.bio || 'No bio yet';
        }
        if (profileAvatarEl && profile.avatar) {
            profileAvatarEl.src = profile.avatar;
        }
        
        if (profileUserId === currentUserId) {
            const firstNameEl = document.getElementById('first_name');
            const lastNameEl = document.getElementById('last_name');
            const bioEl = document.getElementById('bio');
            const avatarUrlEl = document.getElementById('avatar-url');
            
            if (firstNameEl) firstNameEl.value = profile.first_name || '';
            if (lastNameEl) lastNameEl.value = profile.last_name || '';
            if (bioEl) bioEl.value = profile.bio || '';
            if (avatarUrlEl) avatarUrlEl.value = profile.avatar || '';
        }
        
        await loadFollowStats(profile.id);
        
    } catch (error) {
        alert(`Error loading profile: ${error.message}`);
    }
}

async function loadFollowStats(userId) {
    try {
        const followersResponse = await fetch(`/api/follows/followers/${userId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        let followersData = { data: [] };
        if (followersResponse.ok) {
            const followersText = await followersResponse.text();
            if (followersText) {
                followersData = JSON.parse(followersText);
            }
        }
        
        const followingResponse = await fetch(`/api/follows/following/${userId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        let followingData = { data: [] };
        if (followingResponse.ok) {
            const followingText = await followingResponse.text();
            if (followingText) {
                followingData = JSON.parse(followingText);
            }
        }
        
        const followersCountEl = document.getElementById('followers-count');
        const followingCountEl = document.getElementById('following-count');
        
        const followersCount = followersData.data ? followersData.data.length : 0;
        const followingCount = followingData.data ? followingData.data.length : 0;
        
        if (followersCountEl) followersCountEl.textContent = followersCount;
        if (followingCountEl) followingCountEl.textContent = followingCount;
        
        if (profileUserId !== currentUserId && currentUserId) {
            const isFollowing = followersData.data && followersData.data.some(follow => 
                follow.user && follow.user.id === currentUserId
            );
            
            if (typeof updateFollowButton === 'function') {
                updateFollowButton(isFollowing);
            }
        }
        
    } catch (error) {
        const followersCountEl = document.getElementById('followers-count');
        const followingCountEl = document.getElementById('following-count');
        
        if (followersCountEl) followersCountEl.textContent = '0';
        if (followingCountEl) followingCountEl.textContent = '0';
    }
}

async function loadUserPosts() {
    const postsContainer = document.getElementById('user-posts');
    if (!postsContainer) {
        return;
    }
    
    postsContainer.innerHTML = '<div class="loading">Loading posts...</div>';
    
    try {
        const response = await fetch(`/api/posts/user/${profileUserId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to load posts');
        }
        
        const text = await response.text();
        const data = text ? JSON.parse(text) : { data: [] };
        const posts = data.data || [];
        
        const postsCountEl = document.getElementById('posts-count');
        if (postsCountEl) {
            postsCountEl.textContent = posts.length;
        }
        
        if (posts.length === 0) {
            postsContainer.innerHTML = '<p>No posts yet.</p>';
            return;
        }
        
        postsContainer.innerHTML = '';
        posts.forEach(post => {
            const postElement = document.createElement('div');
            postElement.className = 'post-card';
            
            const isOwnPost = currentUserId && post.user && post.user.id === currentUserId;
            const isLiked = post.is_liked || false;
            
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
                    <span class="post-time">${dateString}</span>
                </div>
                <div class="post-content">${post.content || 'No content'}</div>
                ${post.image_url ? `<img src="${post.image_url}" alt="Post image" class="post-image">` : ''}
                <div class="post-actions">
                    <div class="post-action ${isLiked ? 'liked' : ''}" onclick="toggleLike(${post.id}, ${isLiked})">
                        ❤ ${post.like_count || 0} Likes
                    </div>
                    ${isOwnPost ? `
                        <div class="post-action" onclick="editPost(${post.id})">✏️ Edit</div>
                        <div class="post-action" onclick="deletePost(${post.id})">🗑️ Delete</div>
                    ` : ''}
                </div>
            `;
            
            postsContainer.appendChild(postElement);
        });
        
    } catch (error) {
        postsContainer.innerHTML = `<p>Error loading posts: ${error.message}</p>`;
    }
}

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
    
    if (password) {
        updateData.password = password;
    }
    
    try {
        const response = await fetch(`/api/users/profile`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(updateData)
        });
        
        if (!response.ok) {
            throw new Error('Failed to update profile');
        }
        
        alert('Profile updated successfully!');
        
        if (currentUser) {
            const updatedUser = {
                ...currentUser,
                first_name: firstName,
                last_name: lastName,
                bio,
                avatar
            };
            localStorage.setItem('user', JSON.stringify(updatedUser));
        }
        
        await loadProfile();
        
        const postsTabBtn = document.querySelector('[data-tab="posts"]');
        if (postsTabBtn) {
            postsTabBtn.click();
        }
        
    } catch (error) {
        alert(`Error updating profile: ${error.message}`);
    }
}

async function editPost(postId) {
    const newContent = prompt('Edit your post:');
    if (newContent === null) return;
    
    try {
        const response = await fetch(`/api/posts/${postId}`, {
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
        
        await loadUserPosts();
        
    } catch (error) {
        alert(error.message);
    }
}

async function deletePost(postId) {
    if (!confirm('Are you sure you want to delete this post?')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/posts/${postId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to delete post');
        }
        
        await loadUserPosts();
        
    } catch (error) {
        alert(error.message);
    }
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
        const postElement = document.querySelector(`.post-card:has([onclick="toggleLike(${postId}, ${isLiked})"])`);
        if (postElement) {
            const likeAction = postElement.querySelector('.post-action');
            
            if (likeAction) {
                if (isLiked) {
                    likeAction.classList.remove('liked');
                    const likeText = likeAction.textContent.trim();
                    const likeCount = parseInt(likeText.match(/\d+/)[0]) || 0;
                    likeAction.textContent = `❤ ${Math.max(0, likeCount - 1)} Likes`;
                    likeAction.setAttribute('onclick', `toggleLike(${postId}, false)`);
                } else {
                    likeAction.classList.add('liked');
                    const likeText = likeAction.textContent.trim();
                    const likeCount = parseInt(likeText.match(/\d+/)[0]) || 0;
                    likeAction.textContent = `❤ ${likeCount + 1} Likes`;
                    likeAction.setAttribute('onclick', `toggleLike(${postId}, true)`);
                }
            }
        }
    } catch (error) {
        console.error('Error toggling like:', error);
        alert(error.message);
    }
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initProfilePage);
} else {
    initProfilePage();
}