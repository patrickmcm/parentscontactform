document.getElementById('loginForm').addEventListener('submit', function(e) {
    e.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const loginButton = document.getElementById('loginButton');
    const errorMessage = document.getElementById('errorMessage');

    // Clear previous errors
    errorMessage.classList.remove('show');
    loginButton.disabled = true;
    loginButton.textContent = 'Signing in...';

    // Basic validation
    if (!email || !password) {
        showError('Please fill in all fields');
        loginButton.disabled = false;
        loginButton.textContent = 'Sign In';
        return;
    }

    // Submit the form (you can replace this with AJAX if needed)
    this.submit();
});

function showError(message) {
    const errorMessage = document.getElementById('errorMessage');
    errorMessage.textContent = message;
    errorMessage.classList.add('show');
}

// Show error message from server if present
const urlParams = new URLSearchParams(window.location.search);
const error = urlParams.get('error');
if (error) {
    showError(decodeURIComponent(error));
}