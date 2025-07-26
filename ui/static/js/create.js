document.addEventListener('DOMContentLoaded', function() {
    const fileInput = document.getElementById('file');
    const textInput = document.getElementById('text');
    const inputs = document.querySelectorAll('input[name="type"]');
    const urlParams = new URLSearchParams(window.location.search);
    const view = urlParams.get('view');

    function toggleTypeVisibility(e) {
	const fileInputContainer = document.getElementById('file-input');
	const textInputContainer = document.getElementById('text-input');

	if (e.target.id === 'file-radio') {
	    fileInputContainer.classList.remove('hidden');
	    textInputContainer.classList.add('hidden');
	    fileInput.required = true;
	    textInput.required = false;
	} else {
	    textInputContainer.classList.remove('hidden');
	    fileInputContainer.classList.add('hidden');
	    textInput.required = true;
	    fileInput.required = false;
	}
    }

    inputs.forEach(input => {
	input.addEventListener('click', toggleTypeVisibility);
    });

    if (view === 'text') {
	document.getElementById('text-radio').click();
    } else {
	document.getElementById('file-radio').click();
    }

    document.getElementById('create-form').addEventListener('submit', function (e) {
	e.preventDefault();

	const type = document.querySelector('input[name="type"]:checked').value;
	if (type === 'text') {
	    fileInput.value = "";
	} else {
	    textInput.value = "";
	}
	
	this.submit();
    });
});
