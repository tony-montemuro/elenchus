document.addEventListener('DOMContentLoaded', function() {
    const inputs = document.querySelectorAll('input[name="status"]');
    const urlParams = new URLSearchParams(window.location.search);
    const view = urlParams.get('view');

    function toggleQuizzesVisibility(e) {
	const publishedList = document.getElementById('published-list');
	const unpublishedList = document.getElementById('unpublished-list');

	if (e.target.id === 'published') {
	    publishedList.classList.remove('hidden');
	    unpublishedList.classList.add('hidden');
	} else {
	    unpublishedList.classList.remove('hidden');
	    publishedList.classList.add('hidden');
	}
    }

    inputs.forEach(input => {
	input.addEventListener('click', toggleQuizzesVisibility);
    });

    if (view === 'unpublished') {
	document.getElementById('unpublished').click();
    } else {
	document.getElementById('published').click();
    }
});
