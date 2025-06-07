// Check if API_URL is already defined to avoid redeclaration errors
if (typeof API_URL === 'undefined') {
    var API_URL = '/api';
}

// Check if these variables are already defined
if (typeof token === 'undefined') {
    var token = localStorage.getItem('token');
}
if (typeof user === 'undefined') {
    var user = JSON.parse(localStorage.getItem('user'));
}

console.log('🔧 DEBUG: Posts.js loaded');
console.log('Token exists:', !!token);
console.log('User data:', user);

// Initialize posts page
async function initPostsPage() {
    console.log('🔧 DEBUG: Initializing posts page');
    
    // Check if user is logged in
    if (!token) {
        console.log('🔧 DEBUG: No token found, redirecting to login');
        alert('Please login to create posts');
        window.location.href = '/login';
        return;
    }
    
    // Set up post form submission
    const postForm = document.getElementById('post-form');
    if (!postForm) {
        console.log('🔧 DEBUG: post-form element not found');
        return;
    }
    
    console.log('🔧 DEBUG: Setting up post form event listener');
    postForm.addEventListener('submit', createPost);
}

// Create a new post
async function createPost(e) {
    console.log('🔧 DEBUG: Creating new post');
    e.preventDefault();
    
    const contentElement = document.getElementById('post-content');
    const imageUrlElement = document.getElementById('image-url');
    
    if (!contentElement) {
        console.log('🔧 DEBUG: post-content element not found');
        return;
    }
    
    const content = contentElement.value;
    const imageUrl = imageUrlElement ? imageUrlElement.value : '';
    
    console.log('🔧 DEBUG: Post data:');
    console.log('- Content:', content);
    console.log('- Image URL:', imageUrl);
    
    if (!content.trim()) {
        console.log('🔧 DEBUG: Content is empty');
        alert('Please enter some content for your post');
        return;
    }
    
    try {
        console.log('🔧 DEBUG: Making API request to create post');
        console.log('API URL:', `${API_URL}/posts`);
        
        const requestBody = {
            content: content.trim(),
            image_url: imageUrl.trim() || null
        };
        
        console.log('🔧 DEBUG: Request body:', requestBody);
        
        const response = await fetch(`${API_URL}/posts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(requestBody)
        });
        
        console.log('🔧 DEBUG: Create post response:');
        console.log('- Status:', response.status);
        console.log('- OK:', response.ok);
        
        const data = await response.json();
        console.log('🔧 DEBUG: Response data:', data);
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to create post');
        }
        
        // Clear form
        console.log('🔧 DEBUG: Clearing form');
        contentElement.value = '';
        if (imageUrlElement) {
            imageUrlElement.value = '';
        }
        
        // Show success message
        console.log('🔧 DEBUG: Post created successfully');
        alert('Post created successfully!');
        
        // Optionally redirect to home page to see the post
        if (confirm('Post created! Do you want to go to the home page to see it?')) {
            console.log('🔧 DEBUG: Redirecting to home page');
            window.location.href = '/';
        }
    } catch (error) {
        console.error('🔧 DEBUG: Error creating post:', error);
        console.error('🔧 DEBUG: Error stack:', error.stack);
        alert(error.message);
    }
}

// Initialize on page load
console.log('🔧 DEBUG: Setting up posts page initialization');
// Wait for DOM to be ready before initializing
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initPostsPage);
} else {
    initPostsPage();
}