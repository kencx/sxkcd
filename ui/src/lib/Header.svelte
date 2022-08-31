<script lang="ts">
	const theme = {
		light: "light",
		dark: "dark",
	}
	export let systemTheme: string = theme.light;
	export let checked: boolean = false;

	let state: string | undefined
	function toggle(systemTheme: string) {
		// handle system set theme
		state = document.documentElement.dataset['theme']
		if (state === undefined) {
			state = systemTheme
		}
		document.documentElement.dataset['theme'] = state === theme.light ? theme.dark : theme.light
		state = document.documentElement.dataset['theme']
	}

</script>


<header class="container header">
<nav>
	<ul>
	</ul>
	<ul class="right-icons">
		<li>
			<label for="theme-switcher" class="sr-only">Switch theme</label>
			<input id="theme-switcher" type="checkbox" role="switch"
				on:change={() => toggle(systemTheme)} bind:checked={checked}
			/>
		</li>
		<li><a target="_blank" href="https://github.com/kencx/sxkcd" alt="source code">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
			fill="none" stroke={state === theme.light ? "black" : "white"} stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-github"><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path></svg>
		</a></li>
	</ul>
</nav>
</header>

<style>
.header {
	padding: 0;
}

.right-icons {
	padding-right: 1rem;
}

.sr-only {
	position: absolute;
	width: 1px;
	height: 1px;
	padding: 0;
	margin: -1px;
	overflow: hidden;
	clip: rect(0,0,0);
	white-space: nowrap;
	border-width: 0;
}
</style>
