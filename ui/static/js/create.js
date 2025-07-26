document.addEventListener('DOMContentLoaded', function() {
    const inputs = document.querySelectorAll('input[name="type"]');

    function toggleTypeVisibility(e) {
	const pdfInput = document.getElementById('pdf-input');
	const textInput = document.getElementById('text-input');

	if (e.target.id === 'pdf') {
	    pdfInput.classList.remove('hidden');
	    textInput.classList.add('hidden');
	} else {
	    textInput.classList.remove('hidden');
	    pdfInput.classList.add('hidden');
	}
    }

    inputs.forEach(input => {
	input.addEventListener('click', toggleTypeVisibility);
    });

    document.getElementById('pdf').click();
});
