// script.js

// Seleciona o círculo
/*const circle = document.querySelector('.circle');

// Variável para armazenar a cor do círculo
let currentColor = 'rgb(255, 0, 0)'; // Cor inicial

// Função para atualizar a posição do círculo e a cor
function updateCircle(event) {
    // Calcula a nova posição do círculo para seguir o cursor
    const x = event.clientX;
    const y = event.clientY;

    // Atualiza a posição do círculo
    circle.style.left = `${x}px`;
    circle.style.top = `${y}px`;
}

// Função para atualizar a cor do círculo
function updateColor(event) {
    // Calcula a cor com base na posição do mouse
    const x = event.clientX;
    const y = event.clientY;

    const red = Math.round((x / window.innerWidth) * 255);
    const green = Math.round((y / window.innerHeight) * 255);
    const blue = Math.round(((x + y) / (window.innerWidth + window.innerHeight)) * 255);

    // Atualiza a cor atual
    currentColor = `rgb(${red}, ${green}, ${blue})`;
    circle.style.backgroundColor = currentColor;
}

// Atualiza a posição do círculo ao mover o mouse
document.addEventListener('mousemove', updateCircle);

// Atualiza a cor do círculo ao clicar
document.addEventListener('click', updateColor);
