<script lang="ts">
	import Header from '$lib/Header.svelte';
	import Footer from '$lib/Footer.svelte';
	import '../app.css';
</script>

<svelte:head>
	<title>rkcd</title>
	<meta name="description" content="A xkcd search engine" />
	<link rel="stylesheet" href="https://unpkg.com/@picocss/pico@latest/css/pico.min.css">
</svelte:head>

<Header/>
<main>
	<slot />
</main>
<Footer/>

