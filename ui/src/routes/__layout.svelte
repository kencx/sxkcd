<script lang="ts">
	import { onMount } from "svelte";
	import Header from '$lib/Header.svelte';
	import Footer from '$lib/Footer.svelte';
	import '../app.css';

	let systemTheme: string;
	let preferDark: boolean;

	onMount(async() => {
		preferDark = window.matchMedia('(prefers-color-scheme:dark)').matches;
		systemTheme = preferDark ? "dark" : "light";
	});

	$: checked = systemTheme === "light" ? false : true;
</script>

<svelte:head>
	<title>rkcd</title>
	<meta name="description" content="Search for that xkcd you swear you remember" />
	<link rel="stylesheet" href="/css/pico.min.css">
</svelte:head>

<Header {systemTheme} checked={checked}/>
<main>
	<slot />
</main>
<Footer/>

