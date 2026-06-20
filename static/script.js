document.addEventListener('DOMContentLoaded', function () {
  const form = document.querySelector('form');
  form.addEventListener('submit', function () {
    document.getElementById('submit-button').disabled = true;
  });
});
