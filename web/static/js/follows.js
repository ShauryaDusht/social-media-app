// Follow functionality for user interactions

// Initialize follow functionality
function initFollows() {
    // Check if we're on a profile page with a follow button
    const followBtn = document.getElementById('follow-btn');
    if (followBtn) {
        followBtn.addEventListener('click', toggleFollow);
    }
}

// Check if current user is following a profile
async function checkFollowStatus(profileUserId) {
    try {
        const token = localStorage.getItem('token');
        const currentUser = JSON.parse(localStorage.getItem('user'));
        
        if (!token || !currentUser) return false;
        
        // Get followers of the profile user
        const response = await fetch(`/api/follows/followers/${profileUserId}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Failed to check follow status');
        }
        
        const data = await response.json();
        
        // Check if current user is in the followers list
        return data.data.some(follow => follow.user.id === currentUser.id);
    } catch (error) {
        console.error('Error checking follow status:', error);
        return false;
    }
}

// Toggle follow/unfollow
async function toggleFollow() {
    try {
        const token = localStorage.getItem('token');
        const currentUser = JSON.parse(localStorage.getItem('user'));
        
        if (!token || !currentUser) {
            window.location.href = '/login';
            return;
        }
        
        // Get profile user ID from global variable in profile.js or from URL
        let targetUserId;
        if (typeof profileUserId !== 'undefined') {
            targetUserId = profileUserId;
        } else {
            const urlParams = new URLSearchParams(window.location.search);
            if (urlParams.has('id')) {
                targetUserId = parseInt(urlParams.get('id'));
            } else {
                console.error('No target user ID found');
                return;
            }
        }
        
        if (targetUserId === currentUser.id) return; // Can't follow yourself
        
        const followBtn = document.getElementById('follow-btn');
        const isFollowing = followBtn.classList.contains('following');
        
        let url, method, body;
        
        if (isFollowing) {
            // Unfollow request
            url = `/api/follows/${targetUserId}`;
            method = 'DELETE';
            body = null;
        } else {
            // Follow request
            url = `/api/follows`;
            method = 'POST';
            body = JSON.stringify({ user_id: targetUserId });
        }
        
        const response = await fetch(url, {
            method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body
        });
        
        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || `Failed to ${isFollowing ? 'unfollow' : 'follow'} user`);
        }
        
        // Update follow button state
        updateFollowButton(!isFollowing);
        
        // Reload follow stats
        if (typeof loadFollowStats === 'function') {
            loadFollowStats(targetUserId);
        }
    } catch (error) {
        console.error('Error toggling follow:', error);
        alert(error.message);
    }
}

// Update follow button appearance
function updateFollowButton(isFollowing) {
    const followBtn = document.getElementById('follow-btn');
    if (followBtn) {
        if (isFollowing) {
            followBtn.textContent = 'Unfollow';
            followBtn.classList.add('following');
        } else {
            followBtn.textContent = 'Follow';
            followBtn.classList.remove('following');
        }
    }
}

// Run on page load
document.addEventListener('DOMContentLoaded', initFollows);