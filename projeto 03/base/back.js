const proxyUrl = 'http://localhost:3000/check-url';

function verificarSites() {
    listaSites.forEach(url => {
        fetch(`${proxyUrl}?url=${encodeURIComponent(url)}`)
            .then(response => response.text())
            .then(text => {
                const statusElement = document.getElementById(`status-${encodeURIComponent(url)}`);
                if (text.includes('online')) {
                    statusElement.textContent = 'Online';
                    statusElement.className = 'status online';
                } else {
                    statusElement.textContent = 'Offline';
                    statusElement.className = 'status offline';
                }
            })
            .catch(() => {
                const statusElement = document.getElementById(`status-${encodeURIComponent(url)}`);
                statusElement.textContent = 'Offline';
                statusElement.className = 'status offline';
            });
    });
}
