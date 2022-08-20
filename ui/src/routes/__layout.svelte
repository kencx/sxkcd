<script lang="ts">
	import { onMount } from "svelte";
	import Header from '$lib/Header.svelte';
	import Footer from '$lib/Footer.svelte';
	import '../app.css';

	let theme: string;
	let preferDark: boolean;

	onMount(async() => {
		preferDark = window.matchMedia('(prefers-color-scheme:dark)').matches;
		theme = preferDark ? "dark" : "light";
	});

	$: checked = theme === "light" ? false : true;
</script>

<svelte:head>
	<title>rkcd</title>
	<meta name="description" content="A xkcd search engine" />
	<link rel="stylesheet" href="https://unpkg.com/@picocss/pico@latest/css/pico.min.css">
</svelte:head>

<Header currentTheme={theme} checked={checked}/>
<main>
	<slot />
</main>
<Footer/>

