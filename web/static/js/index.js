// Base API URL
const API_URL = '/api';

// Initialize home page
function initHomePage() {
    const token = localStorage.getItem('token');
    const navMenu = document.getElementById('nav-menu');
    const heroButtons = document.getElementById('hero-buttons');
    const timelineContainer = document.getElementById('timeline-container');
    
    if (token) {
        // User is logged in
        navMenu.innerHTML = `
            <li class="nav-item"><a href="/web/index.html" class="nav-link active">Home</a></li>
            <li class="nav-item"><a href="/web/posts.html" class="nav-link">Posts</a></li>
            <li class="nav-item"><a href="/web/profile.html" class="nav-link">Profile</a></li>
            <li class="nav-item"><a href="#" id="logout-link" class="nav-link">Logout</a></li>
        `;
        
        heroButtons.innerHTML = `
            <a href="/web/posts.html" class="btn btn-primary">View Posts</a>
            <a href="/web/profile.html" class="btn">My Profile</a>
        `;
        
        // Add logout functionality
        document.getElementById('logout-link').addEventListener('click', (e) => {
            e.preventDefault();
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.reload();
        });
        
        // Load timeline
        loadTimeline(timelineContainer);
    } else {
        // User is not logged in
        navMenu.innerHTML = `
            <li class="nav-item"><a href="/web/index.html" class="nav-link active">Home</a></li>
            <li class="nav-item"><a href="/web/login.html" class="nav-link">Login</a></li>
            <li class="nav-item"><a href="/web/signup.html" class="nav-link">Sign Up</a></li>
        `;
        
        heroButtons.innerHTML = `
            <a href="/web/login.html" class="btn btn-primary">Login</a>
            <a href="/web/signup.html" class="btn">Sign Up</a>
        `;
    }
}

// Load timeline for logged-in users
async function loadTimeline(container) {
    const token = localStorage.getItem('token');
    if (!token) return;
    
    container.innerHTML = `
        <h2>Your Timeline</h2>
        <div id="timeline-posts" class="posts-list">
            <div class="loading">Loading timeline...</div>
        </div>
    `;
    
    try {
        const response = await fetch(`${API_URL}/timeline`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to load timeline');
        }
        
        const timelinePosts = document.getElementById('timeline-posts');
        
        if (data.data.length === 0) {
            timelinePosts.innerHTML = '<p>Your timeline is empty. Follow some users to see their posts!</p>';
            return;
        }
        
        // Render timeline posts
        timelinePosts.innerHTML = '';
        data.data.forEach(post => {
            const postElement = document.createElement('div');
            postElement.className = 'post-card';
            
            postElement.innerHTML = `
                <div class="post-header">
                    <img src="${post.user.avatar || '/static/img/default-avatar.png'}" alt="${post.user.username}" class="post-avatar">
                    <span class="post-user">${post.user.first_name} ${post.user.last_name}</span>
                    <span class="post-time">${new Date(post.created_at).toLocaleString()}</span>
                </div>
                <div class="post-content">${post.content}</div>
                ${post.image_url ? `<img src="${post.image_url}" alt="Post image" class="post-image">` : ''}
                <div class="post-actions">
                    <div class="post-action ${post.is_liked ? 'liked' : ''}">
                        ❤ ${post.like_count} Likes
                    </div>
                </div>
            `;
            
            timelinePosts.appendChild(postElement);
        });
    } catch (error) {
        console.error('Error loading timeline:', error);
        container.innerHTML = `<p>Error loading timeline: ${error.message}</p>`;
    }
}

// Initialize on page load
initHomePage();