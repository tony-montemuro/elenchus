document.addEventListener('DOMContentLoaded', function() {
    const closeFlashBtn = document.getElementById("close-flash");

    if (closeFlashBtn) {
	closeFlashBtn.addEventListener("click", (e) => {
	    closeFlashBtn.parentElement.style.display = 'none';
	});
    }
});
