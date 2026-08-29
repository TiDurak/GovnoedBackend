function closeSplash() {
    const splash = document.getElementById('splash');

    splash.classList.add('splash--hidden');

    setTimeout(() => {
        splash.remove();
    }, 300);
}