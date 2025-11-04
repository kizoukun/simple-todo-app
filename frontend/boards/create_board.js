const BASE_API_URL = "http://localhost:8080/api/v1";

async function handleCreateBoard(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const title = formData.get('title');
    const description = formData.get('description');

    try {
        const response = await fetch(`${BASE_API_URL}/boards`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            },
            body: JSON.stringify({ title: title, description: description })
        });
        const data = await response.json();
        if (!data.success) {
            alert("Failed to create board: " + data.message);
            return;
        }
        window.location.href = `./index.html`;
    } catch (error) {
        console.error("Error creating board:", error);
    }
}