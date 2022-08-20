<script context="module" lang="ts">
	export const prerender = true;
</script>

<script lang="ts">
	import Search from "$lib/Search.svelte";
</script>

<div class="container center">
	<div class="title">
		<h1><a class="contrast" href="/"><span>r</span>kcd</a></h1>
	</div>

	<div class="desc">
		<p>Blazing fast <a class="contrast" href="https://xkcd.com">xkcd</a> search engine</p>
	</div>

	<Search/>
</div>

<style>
.center {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	padding-left: 2rem;
	padding-right: 2rem;
}

.title {
	text-align: center;
	margin-bottom: 0;
}

.title h1 {
	font-size: 3.5rem;
	margin-bottom: 0;
}

.title h1 a {
	color: inherit;
	text-decoration: none;
}

.title h1 a:hover {
	text-decoration: underline;
}

.desc {
	/* margin: 1rem; */
	/* line-height: 1rem; */
}
</style>
