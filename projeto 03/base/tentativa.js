const listaSites = [];

// Função para adicionar um site à lista
function adicionarSite() {
    const urlInput = document.getElementById('url');
    let url = urlInput.value.trim();

    if (!url) {
        alert('Por favor, insira uma URL.');
        return;
    }

    // Adicionar http:// se a URL não começar com http:// ou https://
    if (!/^https?:\/\//i.test(url)) {
        url = 'http://' + url;
    }

    // Verifica se a URL já está na lista
    if (listaSites.includes(url)) {
        alert('O site já está na lista.');
        return;
    }

    // Adiciona a URL à lista e atualiza a exibição
    listaSites.push(url);
    atualizarLista();
    urlInput.value = ''; // Limpa o campo de entrada
}

// Função para atualizar a lista exibida
function atualizarLista() {
    const listaElement = document.getElementById('lista-sites');
    listaElement.innerHTML = '';

    listaSites.forEach(url => {
        const li = document.createElement('li');
        li.innerHTML = `${url} - <span class="status" id="status-${encodeURIComponent(url)}">Aguardando...</span>`;
        listaElement.appendChild(li);
    });
}

// Função para verificar o status de todos os sites na lista
function verificarSites() {
    listaSites.forEach(url => {
        fetch(url, { method: 'HEAD' })
            .then(response => {
                const statusElement = document.getElementById(`status-${encodeURIComponent(url)}`);

                console.log(response)
                if (response.ok) {
                    statusElement.textContent = 'Online';
                    statusElement.className = 'status online';
                } else {
                    statusElement.textContent = 'Offline';
                    statusElement.className = 'status offline';
                }
            })
            .catch(() => {
                const statusElement = document.getElementById(`status-${encodeURIComponent(url)}`);
                statusElement.textContent = 'Erro';
                statusElement.className = 'Erro de leitura';
            });
    });
}
