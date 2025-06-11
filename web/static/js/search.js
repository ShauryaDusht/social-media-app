// Search functionality for user search

// Initialize search functionality
function initSearch() {
    const searchInput = document.getElementById('search-input');
    const searchBtn = document.getElementById('search-btn');
    const searchResults = document.getElementById('search-results');

    // Add event listeners
    if (searchBtn) {
        searchBtn.addEventListener('click', performSearch);
    }

    if (searchInput) {
        searchInput.addEventListener('keyup', (e) => {
            if (e.key === 'Enter') {
                performSearch();
            }
        });

        // Clear results when input is cleared
        searchInput.addEventListener('input', () => {
            if (searchInput.value.trim() === '') {
                searchResults.innerHTML = '';
                searchResults.classList.remove('active');
            }
        });
    }

    // Close search results when clicking outside
    document.addEventListener('click', (e) => {
        if (!e.target.closest('.search-container')) {
            searchResults.classList.remove('active');
        }
    });
}

// Perform search
async function performSearch() {
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');
    const query = searchInput.value.trim();

    if (!query) return;

    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/users/search?q=${encodeURIComponent(query)}&limit=10&offset=0`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Search failed');
        }

        const data = await response.json();
        displaySearchResults(data.data || []);
    } catch (error) {
        console.error('Error searching users:', error);
    }
}

// Display search results
function displaySearchResults(users) {
    const searchResults = document.getElementById('search-results');
    searchResults.innerHTML = '';

    if (users.length === 0) {
        searchResults.innerHTML = '<div class="search-result-item">No users found</div>';
        searchResults.classList.add('active');
        return;
    }

    users.forEach(user => {
        const resultItem = document.createElement('div');
        resultItem.className = 'search-result-item';
        resultItem.innerHTML = `
            <img src="${user.avatar || '/static/images/default-avatar.png'}" alt="${user.username}" class="search-result-avatar">
            <span class="search-result-name">${user.first_name} ${user.last_name}</span>
            <span class="search-result-username">@${user.username}</span>
        `;

        resultItem.addEventListener('click', () => {
            window.location.href = `/profile?id=${user.id}`;
        });

        searchResults.appendChild(resultItem);
    });

    searchResults.classList.add('active');
}

// Run on page load
document.addEventListener('DOMContentLoaded', initSearch);