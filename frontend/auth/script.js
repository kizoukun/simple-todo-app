// Tab switching
const tabBtns = document.querySelectorAll('.tab-btn');
const forms = document.querySelectorAll('.form-group');

tabBtns.forEach(btn => {
    btn.addEventListener('click', () => {
        const tabName = btn.dataset.tab;

        // Remove active class from all tabs and forms
        tabBtns.forEach(b => b.classList.remove('active'));
        forms.forEach(f => f.classList.remove('active'));

        // Add active class to clicked tab and corresponding form
        btn.classList.add('active');
        document.querySelector(`.form-group[data-tab="${tabName}"]`).classList.add('active');

        // Clear error messages
        clearAllErrors();
    });
});

// Validation functions
function validateEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
}

function validatePassword(password) {
    return password.length >= 6;
}

function clearAllErrors() {
    document.querySelectorAll('.error-message').forEach(msg => msg.classList.remove('show'));
    document.querySelectorAll('input').forEach(input => input.classList.remove('error'));
}

function showError(inputId, errorId, message) {
    const input = document.getElementById(inputId);
    const errorMsg = document.getElementById(errorId);
    input.classList.add('error');
    errorMsg.textContent = message;
    errorMsg.classList.add('show');
}

// Login form submission
document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    clearAllErrors();

    const email = document.getElementById('loginEmail').value.trim();
    const password = document.getElementById('loginPassword').value;



    const successMsg = document.getElementById('loginSuccess');
    successMsg.textContent = ``;
    successMsg.classList.remove('show');

    let isValid = true;

    if (!email) {
        showError('loginEmail', 'loginEmailError', 'Email is required');
        return
    } else if (!validateEmail(email)) {
        showError('loginEmail', 'loginEmailError', 'Please enter a valid email');
        return;
    }

    if (!password) {
        showError('loginPassword', 'loginPasswordError', 'Password is required');
        return;
    } else if (!validatePassword(password)) {
        showError('loginPassword', 'loginPasswordError', 'Password must be at least 6 characters');
        return;
    }


    const resp = await loginUser(email, password);
    if(!resp.success){
        showError('loginEmail', 'loginEmailError', resp.message || 'Login failed');
        return;
    }

    successMsg.textContent = `Welcome back! Logged in as ${email}`;
    successMsg.classList.add('show');
    document.getElementById('loginForm').reset();
    window.localStorage.setItem('authToken', resp.data.token || "");
    window.location.href = "../boards/index.html";
});

// Register form submission
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    clearAllErrors();

    const firstName = document.getElementById('firstName').value.trim();
    const lastName = document.getElementById('lastName').value.trim();
    const email = document.getElementById('registerEmail').value.trim();
    const password = document.getElementById('registerPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const successMsg = document.getElementById('registerSuccess');
    successMsg.textContent = ``;
    successMsg.classList.remove('show');


    if (!firstName) {
        showError('firstName', 'firstNameError', 'First name is required');
        return;
    }

    if (!lastName) {
        showError('lastName', 'lastNameError', 'Last name is required');
        return;
    }

    if (!email) {
        showError('registerEmail', 'registerEmailError', 'Email is required');
        return
    } else if (!validateEmail(email)) {
        showError('registerEmail', 'registerEmailError', 'Please enter a valid email');
        return;
    }

    if (!password) {
        showError('registerPassword', 'registerPasswordError', 'Password is required');
        return;

    } else if (!validatePassword(password)) {
        showError('registerPassword', 'registerPasswordError', 'Password must be at least 6 characters');
        return;
    }

    if (!confirmPassword) {
        showError('confirmPassword', 'confirmPasswordError', 'Please confirm your password');
        return;
    } else if (password !== confirmPassword) {
        showError('confirmPassword', 'confirmPasswordError', 'Passwords do not match');
        return;
    }

    const resp = await createUser(firstName, lastName, email, password, confirmPassword);
    if(!resp.success){
        showError('registerEmail', 'registerEmailError', resp.message || 'Registration failed');
        return;
    }

    successMsg.textContent = `Account created successfully for ${firstName} ${lastName}!`;
    successMsg.classList.add('show');
    document.getElementById('registerForm').reset();
});

async function createUser(firstName, lastName, email, password, confirm) {
    try {
        const response = await fetch('http://localhost:8080/api/v1/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                first_name: firstName,
                last_name: lastName,
                email: email,
                password: password,
                confirm_password: confirm,
            }),
        });
        const data = await response.json();
        return data;
    } catch (error) {
        console.error("Error creating user:", error);
    }
}

async function loginUser(email, password) {
    try {
        const response = await fetch('http://localhost:8080/api/v1/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: email,
                password: password,
            }),
        });
        const data = await response.json();
        return data;
    } catch (error) {
        console.error("Error logging in user:", error);
    }   
}