<script lang="ts">
	import Comic from '$lib/Comic.svelte';
	import ExampleTable from './ExampleTable.svelte';
	import Select from './Select.svelte';
	import { search } from './search'

	let timer: any;

	const debounce = (e: any) => {
		clearTimeout(timer);
		timer = setTimeout(() => {
			query = e.target.value;
		}, 300);
	}

	function focus() {
		let searchBar: any = document.getElementById("search-bar")
		if (searchBar != null) {
			searchBar.focus();
		}
	}

	function handleKeyDown(event: any) {
		switch (event.key) {
			case "/":
				event.preventDefault()
				focus();
				break;
		}
	}

	const PAGE_SIZE = 20;
	let query: string = '';
	let sorted = 'relevancy';
	let currentItems = PAGE_SIZE;

	function sortComics(comics: Comic[], sorted: string) {
		if (comics.length) {

			if (sorted == 'oldest') {
				comics = comics.sort(function(a,b) {
					return a.num - b.num
				});
			} else if (sorted == 'newest') {
				comics = comics.sort(function(a,b) {
					return b.num - a.num
				});
			} else if (query.startsWith("@date:") || query.startsWith("#")) {
				comics = comics.sort(function(a,b) {
					return a.num - b.num
				});
			} else {
			// original sorting
				comics = comics.sort(function(a,b) {
					return a.id - b.id
				});
			}
		}
		return comics;
	}

	$: promise = search(query);
	$: comics = promise
			.then(result => (result != null) ? sortComics(result.comics, sorted) : null)
</script>

<svelte:window on:keydown={handleKeyDown}/>

<div class="container search">
	<ExampleTable/>

	<label for="search-bar" class="sr-only">Search</label>
	<input id="search-bar" type="search"
		autocomplete="off"
		placeholder="Search..."
		on:change={() => {sorted = 'relevancy'; currentItems = PAGE_SIZE}}
		on:input|preventDefault={debounce}
	/>

	{#await promise then result}
		{#if result != null}

			<div class="options">
				<div class="search-message">
					<p>found <span class="contrast">{result.count}</span> results in {(result.time*1000).toFixed(3)}ms</p>
				</div>
				<Select bind:selected={sorted}/>
			</div>

			<!-- render comics -->
			<div>
			{#await comics then data}
				{#if data != null}

					{#each data.slice(0, currentItems) as comic}
						<Comic result={comic}/>
					{/each}

					<!-- load more button -->
					{#if currentItems < data.length}
						<button class="secondary outline load-more" on:click={() => currentItems += PAGE_SIZE}>
							Show More
						</button>
					{/if}
				{/if}
			{/await}
			</div>
		{/if}

	{:catch error}
		<small class="error">An error occured: {error.message}</small>
	{/await}
</div>

<style>
.search {
	width: 75%;
}

.options {
	display: flex;
	justify-content: space-between;
}

.load-more {
	font-size: 0.8rem;
	width: 127px;
	margin-left: auto;
	margin-right: 0;
}

.search-message {
	margin-top: 0.65rem;
	margin-left: 0.5rem;
}

.search-message p {
	font-size: 0.8rem;
}

.contrast {
	color: orange;
	font-weight: bold;
}

.error {
	display: flex;
	justify-content: center;
	align-items: center;
	margin-top: 0.5rem;
	margin: auto;
	font-size: 0.8rem;
	color: red;
}

@media (min-width:320px) and (max-width:640px) {
	.options {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.load-more {
		margin: 0 auto;
	}
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
