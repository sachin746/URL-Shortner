document.addEventListener('DOMContentLoaded', async () => {
    if (!api.isAuthenticated()) {
        window.location.href = 'login.html';
        return;
    }

    const loading = document.getElementById('loading');
    const content = document.getElementById('profileContent');

    try {
        const response = await api.getUser();
        const user = response.data; // Assuming response structure { data: user, message: ... } based on getuser.go

        document.getElementById('userName').textContent = user.name;
        document.getElementById('userUsername').textContent = `@${user.username}`;
        document.getElementById('userEmail').textContent = user.email;
        document.getElementById('userMobile').textContent = user.mobile;

        // Set avatar initials
        const initials = user.name.split(' ').map(n => n[0]).join('').toUpperCase().substring(0, 2);
        document.getElementById('avatar').textContent = initials;

        loading.classList.add('hidden');
        content.classList.remove('hidden');
    } catch (error) {
        console.error('Failed to load profile:', error);
        alert('Failed to load profile. Please login again.');
        api.clearToken();
        window.location.href = 'login.html';
    }
});

document.getElementById('logoutBtn').addEventListener('click', () => {
    api.clearToken();
    window.location.href = 'login.html';
});
