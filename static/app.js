document.addEventListener('DOMContentLoaded', function () {
    const form = document.querySelector('form');
    if (!form) {
        return;
    }
    form.addEventListener('submit', function () {
        const button = form.querySelector('button[type=submit]');
        if (button) {
            button.textContent = 'Generating...';
        }
    });
});
