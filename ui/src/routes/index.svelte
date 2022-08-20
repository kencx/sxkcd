<script context="module" lang="ts">
	export const prerender = true;
</script>

<script lang="ts">
	import Search from "$lib/Search.svelte";
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<div class="flex flex-col justify-center items-center w-full md:w-1/2 m-auto mt-10 px-8">
	<div class="m-2 leading-6">
		<h1 class="font-bold text-7xl"><span class="text-orange-400">r</span>kcd</h1>
		<h2 class="text-lg text-center">[redis]kcd</h2>
	</div>

	<div class="m-4 leading-4">
		<p class="text-center">Blazing fast <a href="https://xkcd.com">xkcd</a> search engine</p>
	</div>

	<Search/>
</div>

<style>
	a {
		background-color: orange;
		/* background-color: #96a8c8 */
	}
</style>
