// Check auth state on load
document.addEventListener('DOMContentLoaded', () => {
    const navLinks = document.getElementById('navLinks');
    if (navLinks && api.isAuthenticated()) {
        navLinks.innerHTML = `
            <a href="stats.html">Stats</a>
            <a href="tech.html">Tech Stack</a>
            <a href="dashboard.html" class="btn btn-primary">Dashboard</a>
        `;
    }
});

function switchTab(tab) {
    const standardForm = document.getElementById('standardForm');
    const customForm = document.getElementById('customForm');
    const buttons = document.querySelectorAll('.tab-btn');
    const alertBox = document.getElementById('alertBox');
    const resultCard = document.getElementById('resultCard');

    // Reset state
    alertBox.className = 'hidden';
    resultCard.classList.remove('active');

    if (tab === 'standard') {
        standardForm.classList.remove('hidden');
        customForm.classList.add('hidden');
        buttons[0].classList.add('active');
        buttons[1].classList.remove('active');
    } else {
        standardForm.classList.add('hidden');
        customForm.classList.remove('hidden');
        buttons[0].classList.remove('active');
        buttons[1].classList.add('active');
    }
}

function showAlert(message, type = 'error') {
    const alertBox = document.getElementById('alertBox');
    alertBox.className = `alert alert-${type}`;
    alertBox.textContent = message;
    alertBox.classList.remove('hidden');
}

function showResult(shortCode) {
    const resultCard = document.getElementById('resultCard');
    const shortUrlDisplay = document.getElementById('shortUrlDisplay');
    const fullUrl = `${window.location.origin}/${shortCode}`;

    shortUrlDisplay.textContent = fullUrl;
    resultCard.classList.add('active');
    document.getElementById('alertBox').classList.add('hidden');
}

async function handleShorten(e, type) {
    e.preventDefault();
    const btn = e.target.querySelector('button[type="submit"]');
    const originalText = btn.textContent;
    btn.textContent = 'Shortening...';
    btn.disabled = true;

    try {
        let response;
        if (type === 'standard') {
            let url = document.getElementById('stdOriginalUrl').value;
            if (!url.match(/^https?:\/\//i)) {
                url = 'https://' + url;
            }
            const validityPeriod = document.getElementById('stdValidityPeriod').value;
            response = await api.shortenUrl(url, validityPeriod);
        } else {
            let url = document.getElementById('custOriginalUrl').value;
            if (!url.match(/^https?:\/\//i)) {
                url = 'https://' + url;
            }
            const code = document.getElementById('custShortCode').value;
            const validityPeriod = document.getElementById('custValidityPeriod').value;
            response = await api.customShortenUrl(url, code, validityPeriod);
        }

        showResult(response.short_code);
        e.target.reset();
    } catch (error) {
        showAlert(error.message);
    } finally {
        btn.textContent = originalText;
        btn.disabled = false;
    }
}

document.getElementById('standardForm').addEventListener('submit', (e) => handleShorten(e, 'standard'));
document.getElementById('customForm').addEventListener('submit', (e) => handleShorten(e, 'custom'));

function copyToClipboard() {
    const text = document.getElementById('shortUrlDisplay').textContent;
    navigator.clipboard.writeText(text).then(() => {
        const btn = document.querySelector('.result-card button');
        const originalText = btn.textContent;
        btn.textContent = 'Copied!';
        setTimeout(() => btn.textContent = originalText, 2000);
    });
}
