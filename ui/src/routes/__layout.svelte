<script lang="ts">
	import Header from '$lib/Header.svelte';
	import Footer from '$lib/Footer.svelte';
	import '../app.css';
</script>

<main>
	<Header/>
	<slot />
	<Footer/>
</main>

<style>
	a {
		/* background-color: orange; */
		background-color: #96a8c8
	}
</style>
