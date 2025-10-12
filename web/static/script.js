document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.querySelector('.search-input');
    const searchForm = document.querySelector('.search-form');

    if (searchInput && !searchInput.value) {
        searchInput.focus();
    }

    searchForm.addEventListener('submit', function () {
        const button = document.querySelector('.search-button')
        button.textContent = 'searching...'
        button.disabled = true;
    });

    document.addEventListener('keydown', function (e) {
        if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
            e.preventDefault()

            searchInput.focus();
            searchInput.select();
        }
    });

    const query = new URLSearchParams(window.location.search).get('q')
    if (query) {
        const terms = query.toLowerCase().split(/\s+/);
        const snippets = document.querySelectorAll('.result-snippet')

        snippets.forEach(snippet => {
            let html= snippet.innerHTML
        })

    }

});

const style = document.createElement('style');
style.textContent = `
    mark {
        background-color: #fff3cd;
        padding: 1px 2px;
        border-radius: 2px;
    }
`;
document.head.appendChild(style);