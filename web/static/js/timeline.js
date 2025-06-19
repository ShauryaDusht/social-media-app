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
    
    // Check if we have cached timeline data and it's not expired
    const cachedTimeline = getTimelineFromCache();
    if (cachedTimeline) {
        // Display cached timeline
        displayPosts(cachedTimeline, postsContainer);
        
        // If cache is older than 2 minutes, refresh in background
        const cacheTime = localStorage.getItem('timeline_cache_time');
        const cacheAge = Date.now() - parseInt(cacheTime || 0);
        if (cacheAge > 2 * 60 * 1000) { // 2 minutes in milliseconds
            refreshTimelineInBackground();
        }
        return;
    }
    
    // No cache or expired cache, show loading and fetch from server
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
        
        // Cache the timeline data
        saveTimelineToCache(posts);
        
        displayPosts(posts, postsContainer);
        
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

// Save timeline data to localStorage cache
function saveTimelineToCache(posts) {
    try {
        localStorage.setItem('timeline_cache', JSON.stringify(posts));
        localStorage.setItem('timeline_cache_time', Date.now().toString());
    } catch (error) {
        console.error('Error saving timeline to cache:', error);
        // If localStorage is full, clear it and try again
        if (error instanceof DOMException && error.name === 'QuotaExceededError') {
            localStorage.clear();
            try {
                localStorage.setItem('timeline_cache', JSON.stringify(posts));
                localStorage.setItem('timeline_cache_time', Date.now().toString());
            } catch (e) {
                console.error('Still unable to cache timeline after clearing localStorage:', e);
            }
        }
    }
}

// Get timeline data from localStorage cache
function getTimelineFromCache() {
    const cachedData = localStorage.getItem('timeline_cache');
    if (!cachedData) return null;
    
    try {
        return JSON.parse(cachedData);
    } catch (error) {
        console.error('Error parsing cached timeline:', error);
        return null;
    }
}

// Refresh timeline in background without showing loading indicator
async function refreshTimelineInBackground() {
    try {
        const response = await fetch(`${API_URL}/posts`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) return;
        
        const data = await response.json();
        const posts = data.data || data || [];
        
        // Update cache
        saveTimelineToCache(posts);
        
        // Update UI if user is still on the page
        const postsContainer = document.getElementById('timeline-posts');
        if (postsContainer) {
            displayPosts(posts, postsContainer);
        }
    } catch (error) {
        console.error('Background refresh error:', error);
    }
}

// Display posts in the container
function displayPosts(posts, container) {
    if (posts.length === 0) {
        container.innerHTML = `
            <div style="text-align: center; padding: 40px; background: #f8f9fa; border-radius: 8px; margin: 20px 0;">
                <h3>Your Timeline is Empty</h3>
                <p>Follow some users or create your first post to see content here.</p>
                <a href="/posts" style="display: inline-block; background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 10px;">Create Post</a>
                <a href="/users" style="display: inline-block; background: #6c757d; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 10px;">Find Users</a>
            </div>
        `;
        return;
    }
    
    container.innerHTML = '';
    posts.forEach(post => {
        const postElement = createTimelinePostElement(post);
        container.appendChild(postElement);
    });
}

function createTimelinePostElement(post) {
    const postElement = document.createElement('div');
    postElement.className = 'post-card';
    postElement.dataset.id = post.id;
    
    // Check if current user liked this post
    const isLiked = post.liked_by && post.liked_by.includes(user.id);
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
        
        const responseData = await response.json();
        const data = responseData.data || responseData;
        
        // Update UI immediately with server response
        const postElement = document.querySelector(`.post-card[data-id="${postId}"]`);
        if (postElement) {
            const likeAction = postElement.querySelector('.post-action');
            const likeCount = postElement.querySelector('.like-count');
            
            if (likeAction && likeCount) {
                if (data.is_liked) {
                    likeAction.classList.add('liked');
                    likeAction.setAttribute('onclick', `toggleLike(${postId}, true)`);
                } else {
                    likeAction.classList.remove('liked');
                    likeAction.setAttribute('onclick', `toggleLike(${postId}, false)`);
                }
                likeCount.textContent = data.like_count || 0;
            }
        }
        
        // Update the cached timeline with server response
        updateCachedPostFromServer(postId, data);
        
        // Force refresh timeline in background to ensure consistency
        setTimeout(() => {
            refreshTimelineInBackground();
        }, 500);
        
    } catch (error) {
        console.error('Error toggling like:', error);
        alert(error.message);
    }
}

// Update a single post in the cached timeline with server data
function updateCachedPostFromServer(postId, serverData) {
    const cachedTimeline = getTimelineFromCache();
    if (!cachedTimeline) return;
    
    const updatedTimeline = cachedTimeline.map(post => {
        if (post.id === postId) {
            post.like_count = serverData.like_count || 0;
            post.liked_by = serverData.liked_by || [];
            post.is_liked = serverData.is_liked || false;
        }
        return post;
    });
    
    saveTimelineToCache(updatedTimeline);
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
        
        // Remove from cache
        removePostFromCache(postId);
        
        // Remove from UI
        const postElement = document.querySelector(`.post-card[data-id="${postId}"]`);
        if (postElement) {
            postElement.remove();
        }
        
        // Refresh timeline to ensure consistency
        setTimeout(() => {
            refreshTimelineInBackground();
        }, 500);
        
    } catch (error) {
        alert(error.message);
    }
}

// Remove a post from the cached timeline
function removePostFromCache(postId) {
    const cachedTimeline = getTimelineFromCache();
    if (!cachedTimeline) return;
    
    const updatedTimeline = cachedTimeline.filter(post => post.id !== postId);
    saveTimelineToCache(updatedTimeline);
}

async function editPost(postId) {
    const cachedTimeline = getTimelineFromCache();
    const post = cachedTimeline?.find(p => p.id === postId);
    const currentContent = post?.content || '';
    
    const newContent = prompt('Edit your post:', currentContent);
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
        
        // Update in cache
        updatePostContentInCache(postId, newContent);
        
        // Update in UI
        const postElement = document.querySelector(`.post-card[data-id="${postId}"]`);
        if (postElement) {
            const contentElement = postElement.querySelector('.post-content');
            if (contentElement) {
                contentElement.textContent = newContent;
            }
        }
        
        // Refresh timeline to ensure consistency
        setTimeout(() => {
            refreshTimelineInBackground();
        }, 500);
        
    } catch (error) {
        alert(error.message);
    }
}

// Update post content in the cached timeline
function updatePostContentInCache(postId, newContent) {
    const cachedTimeline = getTimelineFromCache();
    if (!cachedTimeline) return;
    
    const updatedTimeline = cachedTimeline.map(post => {
        if (post.id === postId) {
            post.content = newContent;
        }
        return post;
    });
    
    saveTimelineToCache(updatedTimeline);
}

function viewProfile(userId) {
    if (userId) {
        window.location.href = `/profile?id=${userId}`;
    }
}

// Clear timeline cache when user logs out
document.addEventListener('logout', function() {
    localStorage.removeItem('timeline_cache');
    localStorage.removeItem('timeline_cache_time');
});

// Force refresh timeline every 2 minutes to ensure consistency
setInterval(() => {
    if (token && document.getElementById('timeline-posts')) {
        refreshTimelineInBackground();
    }
}, 2 * 60 * 1000);  // 2 minutes in milliseconds

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initTimelinePage);
} else {
    initTimelinePage();
}