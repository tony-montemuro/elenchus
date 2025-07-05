document.addEventListener('DOMContentLoaded', function() {
    const inputs = document.querySelectorAll('input[name="status"]');

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
    document.getElementById('published').click();
});
