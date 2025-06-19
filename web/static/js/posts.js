// Helper functions for authentication
function getToken() {
    return localStorage.getItem('token');
}

function getUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
}

function checkAuth() {
    const token = getToken();
    const user = getUser();
    
    if (!token || !user) {
        alert('Please login to access this page');
        window.location.href = '/login';
        return false;
    }
    return true;
}

// Initialize posts page
async function initPostsPage() {
    // Check authentication
    if (!checkAuth()) {
        return;
    }

    const token = getToken();
    const user = getUser();
    
    // Set up post form submission
    const postForm = document.getElementById('post-form');
    if (!postForm) {
        return;
    }
    
    postForm.addEventListener('submit', createPost);
}

// Create a new post
async function createPost(e) {
    e.preventDefault();
    
    const token = getToken();
    if (!token) {
        alert('Please login to create posts');
        window.location.href = '/login';
        return;
    }
    
    const contentElement = document.getElementById('post-content');
    const imageUrlElement = document.getElementById('image-url');
    
    if (!contentElement) {
        return;
    }
    
    const content = contentElement.value;
    const imageUrl = imageUrlElement ? imageUrlElement.value : '';
    
    if (!content.trim()) {
        alert('Please enter some content for your post');
        return;
    }
    
    // Disable form to prevent double submission
    const submitButton = e.target.querySelector('button[type="submit"]');
    if (submitButton) {
        submitButton.disabled = true;
        submitButton.textContent = 'Creating...';
    }
    
    try {
        const requestBody = {
            content: content.trim(),
            image_url: imageUrl.trim() || ''
        };
        
        const response = await fetch('/api/posts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(requestBody)
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || data.message || `HTTP ${response.status}: ${response.statusText}`);
        }
        
        // Clear form
        contentElement.value = '';
        if (imageUrlElement) {
            imageUrlElement.value = '';
        }
        
        // Clear timeline cache so the new post appears for all users
        clearTimelineCache();
        
        // Show success message
        alert('Post created successfully!');
        
        // Optionally redirect to home page
        if (confirm('Post created! Do you want to go to the home page to see it?')) {
            window.location.href = '/';
        }
        
    } catch (error) {
        console.error('Error creating post:', error);
        alert(`Error creating post: ${error.message}`);
    } finally {
        // Re-enable form
        if (submitButton) {
            submitButton.disabled = false;
            submitButton.textContent = 'Create Post';
        }
    }
}

// Function to clear timeline cache
function clearTimelineCache() {
    try {
        localStorage.removeItem('timeline_cache');
        localStorage.removeItem('timeline_cache_time');
        console.log('Timeline cache cleared after creating post');
    } catch (error) {
        console.error('Error clearing timeline cache:', error);
    }
}

// Wait for DOM and try to initialize
function tryInit() {
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initPostsPage);
    } else {
        // DOM is already ready, initialize immediately
        initPostsPage();
    }
}

// Call initialization
tryInit();