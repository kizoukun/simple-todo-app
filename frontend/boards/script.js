const BASE_API_URL = "http://localhost:8080/api/v1";


document.addEventListener("DOMContentLoaded", onLoad)

async function onLoad() {
    const response = await getBoards();
    if (!response.success) {
        console.error("Failed to fetch boards:", response.message);
        return;
    }
    const boards = response.data;
    const ownedBoardList = document.getElementById('ownedBoardsList');
    if (ownedBoardList === null) return;
    ownedBoardList.innerHTML = '';
    for (const board of boards.owned_boards) {
        const boardItem = document.createElement('a');
        boardItem.href = `../index.html?id=${board.id}`;
        boardItem.className = 'block';
        boardItem.innerHTML = ` 
             <div class="bg-blue-200 rounded-lg p-2">
             <p>${board.title}</p>
        <p class="text-sm text-gray-600">${board.description}</p>
             </div>
         `;
        ownedBoardList.appendChild(boardItem);
    }

    const sharedBoardsList = document.getElementById('sharedBoardsList');
    if (sharedBoardsList === null) return;
    sharedBoardsList.innerHTML = '';

    for (const board of boards.team_boards) {
        const boardItem = document.createElement('a');
        boardItem.href = `../index.html?id=${board.id}`;
        boardItem.className = 'block';
        boardItem.innerHTML = ` 
             <div class="bg-blue-200 rounded-lg p-2">
             <p class="text-lg font-semibold">${board.title}</p>
             <p class="text-sm text-gray-600">${board.description}</p>
             </div>
         `;
        sharedBoardsList.appendChild(boardItem);
    }

}

async function getBoards() {
    try {
        const response = await fetch(`${BASE_API_URL}/boards`, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            }
        });
        if (response.status === 401) {
            window.location.href = "../auth/auth.html";
            return null;
        }
        const boards = await response.json();
        return boards;
    } catch (error) {
        console.error("Error fetching boards:", error);
    }
}
