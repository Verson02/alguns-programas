const siteList = document.getElementById("siteList");
const addSiteForm = document.getElementById("addSiteForm");

addSiteForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    let siteUrl = document.getElementById("siteUrl").value.trim();
  
    if (!siteUrl.startsWith("http://") && !siteUrl.startsWith("https://")) {
      siteUrl = `https://${siteUrl}`;
    }
  
    try {
      const response = await fetch("/api/add-site", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url: siteUrl }),
      });
  
      if (response.ok) {
        document.getElementById("siteUrl").value = "";
        fetchSiteStatus();
      } else {
        const errorText = await response.text();
        alert(`Erro ao adicionar site: ${errorText}`);
      }
    } catch (err) {
      console.error("Erro de conexão:", err);
      alert("Erro ao conectar com o servidor.");
    }
  });

async function fetchSiteStatus() {
  try {
    const response = await fetch("/api/check-sites");
    if (response.ok) {
      const sites = await response.json();
      siteList.innerHTML = "";
      sites.forEach((site) => {
        const listItem = document.createElement("li");
        listItem.textContent = `${site.url} - ${site.status} (Última verificação: ${new Date(
          site.checked_at
        ).toLocaleTimeString()})`;
        siteList.appendChild(listItem);
      });
    }
  } catch (err) {
    console.error(err);
  }
}

// Atualiza a lista a cada 10 segundos
setInterval(fetchSiteStatus, 10000);
fetchSiteStatus();