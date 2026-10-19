(() => {
  const article = document.querySelector(".markdown-body");

  // <base> points at the markdown file's folder (so relative images work),
  // which would make "#section" links navigate away. Handle them in place.
  document.addEventListener("click", (e) => {
    const a = e.target.closest("a[href^='#']");
    if (!a) return;
    e.preventDefault();
    const id = decodeURIComponent(a.getAttribute("href").slice(1));
    const target = document.getElementById(id);
    if (!target) return;
    // replace() keeps history at one entry, which browsers require before a page may close itself.
    if (location.hash === "#" + id) target.scrollIntoView();
    else location.replace(location.href.split("#")[0] + "#" + id);
  });

  // Close button. If the browser refuses (the tab has history), say how to close it instead.
  const closeButton = document.querySelector(".close-page");
  const closeIcon = closeButton.innerHTML;
  closeButton.addEventListener("click", () => {
    window.close();
    setTimeout(() => {
      if (window.closed) return;
      const shortcut = /Mac|iPhone|iPad/.test(navigator.platform) ? "⌘W" : "Ctrl+W";
      closeButton.textContent = "Press " + shortcut + " to close";
      setTimeout(() => (closeButton.innerHTML = closeIcon), 2500);
    }, 150);
  });

  // Hover anchors on headings.
  for (const h of article.querySelectorAll("h1[id], h2[id], h3[id], h4[id], h5[id], h6[id]")) {
    const a = document.createElement("a");
    a.className = "heading-anchor";
    a.setAttribute("href", "#" + h.id);
    a.setAttribute("aria-label", "Link to this section");
    a.textContent = "#";
    h.prepend(a);
  }

  // Copy buttons on code blocks.
  document.addEventListener("click", async (e) => {
    const button = e.target.closest(".code-copy");
    if (!button) return;
    const code = button.parentElement.querySelector("pre").textContent;
    try {
      await navigator.clipboard.writeText(code);
    } catch {
      const area = document.createElement("textarea");
      area.value = code;
      document.body.append(area);
      area.select();
      document.execCommand("copy");
      area.remove();
    }
    button.textContent = "Copied";
    setTimeout(() => (button.textContent = "Copy"), 1500);
  });

  // Click a rendered diagram to toggle between fit-to-width and full size.
  document.addEventListener("click", (e) => {
    const box = e.target.closest(".diagram.rendered");
    if (box) box.classList.toggle("full");
  });

  // Reading progress and current-section highlight in the table of contents.
  const progress = document.querySelector(".progress");
  const links = [...document.querySelectorAll(".toc a")];
  const headings = links.map((a) => document.getElementById(decodeURIComponent(a.getAttribute("href").slice(1))));
  const offset = parseInt(getComputedStyle(document.documentElement).getPropertyValue("scroll-padding-top")) || 60;

  let queued = false;
  const update = () => {
    queued = false;
    const max = document.documentElement.scrollHeight - innerHeight;
    progress.style.width = max > 0 ? (scrollY / max) * 100 + "%" : "0";

    let current = 0;
    headings.forEach((h, i) => {
      if (h && h.getBoundingClientRect().top <= offset + 8) current = i;
    });
    if (max > 0 && scrollY >= max - 2) current = headings.length - 1;
    links.forEach((a, i) => a.classList.toggle("active", i === current));
  };
  addEventListener("scroll", () => {
    if (!queued) {
      queued = true;
      requestAnimationFrame(update);
    }
  }, { passive: true });
  update();
})();
