class FichaRPG {
    constructor() {
        this.config = {
            modificadoresRaciais: {
                "Anão": { constituicao: 2 },
                "Elfo": { destreza: 2 },
                "Halfling": { destreza: 2 },
                "Humano": { forca: 1, destreza: 1, constituicao: 1, inteligencia: 1, sabedoria: 1, carisma: 1 },
                "Draconato": { forca: 2, carisma: 1 },
                "Gnomo": { inteligencia: 2 },
                "Meio-Elfo": { carisma: 2, opcionais: ['+2 em dois atributos'] },
                "Meio-Orc": { forca: 2, constituicao: 1 },
                "Tiefling": { inteligencia: 1, carisma: 2 }
            },

            periciasClasse: {
                "Bárbaro": ["Adestrar Animais", "Atletismo", "Intimidação", "Natureza", "Percepção", "Sobrevivência"],
                "Bardo": ["Acrobacia", "Atuação", "Enganação", "História", "Intuição", "Investigação", "Medicina", "Persuasão"],
                "Bruxo": ["Arcanismo", "Enganação", "Intuição", "Investigação", "Religião"],
                "Clérigo": ["Arcanismo", "História", "Intimidação", "Investigação", "Religião"],
                "Druida": ["Acrobacia", "Atletismo", "História", "Intimidação", "Investigação", "Natureza", "Percepção", "Sobrevivência"],
                "Feiticeiro": ["Arcanismo", "Enganação", "História", "Intuição", "Investigação", "Religião"],
                "Guerreiro": ["Acrobacia", "Atletismo", "História", "Intimidação", "Percepção", "Sobrevivência"],
                "Ladino": ["Acrobacia", "Enganação", "Furtividade", "Intuição", "Investigação", "Percepção", "Persuasão"],
                "Mago": ["Arcanismo", "História", "Intuição", "Investigação", "Religião"],
                "Monge": ["Acrobacia", "Atletismo", "Furtividade", "Intuição", "Percepção", "Sobrevivência"],
                "Paladino": ["Atletismo", "Intimidação", "Medicina", "Persuasão", "Religião"],
                "Patrulheiro": ["Acrobacia", "Sobrevivência", "Natureza", "Percepção"]
            },

            periciasRaciais: {
                "Anão": ["História (Anões)", "Ferramentas de Ferreiro"],
                "Elfo": ["Percepção", "Furtividade"],
                "Halfling": ["Furtividade", "Prestidigitação"],
                "Humano": ["Persuasão", "Atletismo"],
                "Draconato": ["Intimidação", "Percepção"],
                "Gnomo": ["Arcanismo", "Prestidigitação"],
                "Meio-Elfo": ["Persuasão", "Percepção"],
                "Meio-Orc": ["Intimidação", "Atletismo"],
                "Tiefling": ["Intuição", "Furtividade"]
            },

            dadosVida: {
                "Bárbaro": 12, "Bardo": 8, "Bruxo": 8, "Clérigo": 8, 
                "Druida": 8, "Feiticeiro": 6, "Guerreiro": 10, 
                "Ladino": 8, "Mago": 6, "Monge": 8, 
                "Paladino": 10, "Patrulheiro": 10
            },

            limitesPericias: {
                "Bárbaro": 2, "Bardo": 3, "Bruxo": 2, "Clérigo": 2,
                "Druida": 2, "Feiticeiro": 2, "Guerreiro": 2,
                "Ladino": 4, "Mago": 2, "Monge": 2,
                "Paladino": 2, "Patrulheiro": 3
            }
        };

        this.inventario = new Inventario();
        this.inicializarEventos();
    }

    inicializarEventos() {
        document.getElementById('gerarAtributos').addEventListener('click', () => this.gerarAtributosAleatorios());
        document.getElementById('classe').addEventListener('change', (e) => this.carregarPericias(e.target.value));
        document.getElementById('raca').addEventListener('change', () => this.aplicarModificadoresRaciais());
        document.getElementById('fichaForm').addEventListener('submit', (e) => this.gerarFicha(e));
        document.getElementById('exportPDF').addEventListener('click', () => this.gerarPDF());
    }

    gerarAtributosAleatorios() {
        const rolarD6 = () => Math.floor(Math.random() * 6) + 1;
        
        const calcularAtributo = () => {
            const dados = Array.from({length: 4}, rolarD6).sort((a, b) => a - b).slice(1);
            return dados.reduce((a, b) => a + b, 0);
        };

        ['forca', 'destreza', 'constituicao', 'inteligencia', 'sabedoria', 'carisma'].forEach(id => {
            document.getElementById(id).value = calcularAtributo();
        });
    }

    carregarPericias(classeSelecionada) {
        const raca = document.getElementById('raca').value;
        const periciasClasse = this.config.periciasClasse[classeSelecionada] || [];
        const periciasRaca = this.config.periciasRaciais[raca] || [];
        
        // Filtrar perícias da classe que não são raciais
        const periciasDisponiveis = periciasClasse.filter(p => !periciasRaca.includes(p));
        
        const container = document.getElementById('periciasLista');
        container.innerHTML = '';
        
        periciasDisponiveis.forEach(pericia => {
            const div = document.createElement('div');
            div.className = 'col-md-6 mb-2';
            div.innerHTML = `
                <div class="form-check">
                    <input class="form-check-input" type="checkbox" name="pericias" id="pericia-${pericia}">
                    <label class="form-check-label" for="pericia-${pericia}">${pericia}</label>
                </div>
            `;
            container.appendChild(div);
        });

        const limite = this.config.limitesPericias[classeSelecionada] || 0;
        this.configurarValidacaoPericias(limite);
    }

    configurarValidacaoPericias(limite) {
        const checkboxes = document.querySelectorAll('input[name="pericias"]');
        checkboxes.forEach(checkbox => {
            checkbox.addEventListener('change', (e) => {
                const selecionados = document.querySelectorAll('input[name="pericias"]:checked').length;
                if(selecionados > limite) {
                    e.target.checked = false;
                    alert(`Você pode selecionar no máximo ${limite} perícias adicionais!`);
                }
            });
        });
    }

    aplicarModificadoresRaciais() {
        const raca = document.getElementById('raca').value;
        const modificadores = this.config.modificadoresRaciais[raca] || {};
        
        Object.entries(modificadores).forEach(([atributo, valor]) => {
            if(typeof valor === 'number') {
                const input = document.getElementById(atributo);
                if(input) input.value = parseInt(input.value) + valor;
            }
        });
    }

    calcularModificador(valor) {
        return Math.floor((valor - 10) / 2);
    }

    gerarFicha(e) {
        e.preventDefault();
        const dados = this.coletarDados();
        const fichaHTML = this.criarTemplateFicha(dados);
        
        document.getElementById('fichaContainer').style.display = 'block';
        document.getElementById('fichaContent').innerHTML = fichaHTML;
    }

    coletarDados() {
        const raca = document.getElementById('raca').value;
        return {
            nome: document.getElementById('nome').value,
            nivel: parseInt(document.getElementById('nivel').value),
            raca: raca,
            classe: document.getElementById('classe').value,
            atributos: this.obterAtributos(),
            periciasClasse: Array.from(document.querySelectorAll('input[name="pericias"]:checked'))
                              .map(input => input.nextElementSibling.textContent),
            periciasRaciais: this.config.periciasRaciais[raca] || [],
            inventario: this.inventario.obterItens(),
            historia: document.getElementById('historiaEditor').innerHTML
        };
    }

    obterAtributos() {
        return {
            forca: parseInt(document.getElementById('forca').value),
            destreza: parseInt(document.getElementById('destreza').value),
            constituicao: parseInt(document.getElementById('constituicao').value),
            inteligencia: parseInt(document.getElementById('inteligencia').value),
            sabedoria: parseInt(document.getElementById('sabedoria').value),
            carisma: parseInt(document.getElementById('carisma').value)
        };
    }

    criarTemplateFicha(dados) {
        const modificadores = Object.fromEntries(
            Object.entries(dados.atributos).map(([chave, valor]) => 
                [chave, this.calcularModificador(valor)]
            )
        );
        
        const PV = this.calcularPV(dados);
        const CA = 10 + modificadores.destreza;

        return `
            <div class="ficha-header">
                <h2>${dados.nome}</h2>
                <h4>${dados.raca} ${dados.classe} - Nível ${dados.nivel}</h4>
            </div>

            <div class="atributos-section">
                ${Object.entries(dados.atributos).map(([atributo, valor]) => `
                    <div class="atributo">
                        <span class="nome-atributo">${atributo.toUpperCase()}</span>
                        <span class="valor-atributo">${valor}</span>
                        <span class="modificador">(${modificadores[atributo] >= 0 ? '+' : ''}${modificadores[atributo]})</span>
                    </div>
                `).join('')}
            </div>

            <div class="status-section">
                <div class="status-item">
                    <span>Pontos de Vida (PV):</span>
                    <span class="valor-status">${PV}</span>
                </div>
                <div class="status-item">
                    <span>Classe de Armadura (CA):</span>
                    <span class="valor-status">${CA}</span>
                </div>
            </div>

            <div class="pericias-section">
                <h3>Perícias</h3>
                <div class="pericias-lista">
                    ${dados.periciasRaciais.map(p => `
                        <div class="pericia-racial">🛡️ ${p}</div>
                    `).join('')}
                    ${dados.periciasClasse.map(p => `
                        <div class="pericia-classe">✔ ${p}</div>
                    `).join('')}
                </div>
            </div>

            <div class="inventario-section">
                <h3>Inventário</h3>
                ${dados.inventario.map(item => `<div class="item">⚔️ ${item}</div>`).join('')}
            </div>

            <div class="historia-section">
                <h3>História</h3>
                <div class="conteudo-historia">${dados.historia}</div>
            </div>
        `;
    }

    calcularPV(dados) {
        const dadoVida = this.config.dadosVida[dados.classe] || 8;
        const modificador = this.calcularModificador(dados.atributos.constituicao);
        return (dadoVida + modificador) * dados.nivel;
    }

    gerarPDF() {
        const doc = new jspdf.jsPDF();
        doc.html(document.getElementById('fichaContent'), {
            margin: [15, 15, 15, 15],
            width: 180,
            callback: () => doc.save('ficha-dnd.pdf')
        });
    }
}

class Inventario {
    constructor() {
        this.itens = [];
        this.container = document.getElementById('inventarioContainer');
        this.inicializar();
    }

    inicializar() {
        this.configurarDragAndDrop();
        document.getElementById('adicionarItem').addEventListener('click', () => this.adicionarItem());
    }

    adicionarItem() {
        const nomeItem = prompt('Nome do item:');
        if(nomeItem) {
            this.itens.push(nomeItem);
            this.atualizarInventario();
        }
    }

    obterItens() {
        return [...this.itens];
    }

    atualizarInventario() {
        this.container.innerHTML = this.itens
            .map((item, index) => `
                <div class="item-inventario" draggable="true" data-index="${index}">${item}</div>
            `).join('');
    }

    configurarDragAndDrop() {
        this.container.addEventListener('dragover', (e) => this.dragOver(e));
        this.container.addEventListener('drop', (e) => this.drop(e));
    }

    dragOver(e) {
        e.preventDefault();
    }

    drop(e) {
        e.preventDefault();
        const origemIndex = parseInt(e.dataTransfer.getData('text/plain'));
        const destinoIndex = this.calcularPosicao(e.clientY);
        
        const [item] = this.itens.splice(origemIndex, 1);
        this.itens.splice(destinoIndex, 0, item);
        this.atualizarInventario();
    }

    calcularPosicao(y) {
        const itens = [...this.container.children];
        return itens.findIndex(item => {
            const rect = item.getBoundingClientRect();
            return y < rect.top + rect.height/2;
        });
    }
}

new FichaRPG();