document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('unpublish').addEventListener('submit', e => {
	if (!confirm('Unpublishing this quiz will make it inaccessible to the public, and remove all quiz attempts from all users. Are you sure you wish to proceed?')) {
	    e.preventDefault();
	}
    });
});
