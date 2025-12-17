const moon = document.getElementById('moon-icon');
const sun = document.getElementById('sun-icon');

function setTheme(isDark) {
	if (isDark) {
		document.documentElement.classList.add('dark');
		localStorage.setItem('theme', 'dark');
	} else {
		document.documentElement.classList.remove('dark');
		localStorage.setItem('theme', 'light');
	}
	
	if (moon && sun) {
		if (isDark) {
			moon.hidden = true;
			sun.hidden = false;
		} else {
			moon.hidden = false;
			sun.hidden = true;
		}
	}
}

function toggleTheme() {
	const isDark = document.documentElement.classList.contains('dark');
	setTheme(!isDark);
}

// Init theme
const saved = localStorage.getItem('theme');
const prefer = window.matchMedia('(prefers-color-scheme: dark)').matches;

if (saved === 'dark') {
	setTheme(true);
} else if (saved === 'light') {
	setTheme(false);
} else if (prefer) {
	setTheme(true);
} else {
	setTheme(false);
}


 