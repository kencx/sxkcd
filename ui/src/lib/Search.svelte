<script lang="ts">
	import Comic from '$lib/Comic.svelte';
	import { search } from './search'

	let query: string = '';
	let page: number = 1;
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

	$: promise = search(query, page);
</script>

<svelte:window on:keydown={handleKeyDown}/>

<div class="container search">
	<label for="search-bar" class="sr-only">Search</label>
	<input id="search-bar" type="search"
		autocomplete="off"
		placeholder="Search..."
		on:input|preventDefault={debounce}
	/>

	{#await promise then data}
		{#if data != null}
			<div class="search-message">
				<p>found <span class="contrast">{data.count}</span> results in {(data.time*1000).toFixed(3)}ms</p>
			</div>
			<div>
				{#if data.comics}
					{#each data.comics as comic}
						<Comic result={comic}/>
					{/each}
				{/if}
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

.search-message {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	margin-bottom: 1.5rem;
}

.search-message p {
	font-size: 0.8rem;
}

.contrast {
	color: orange;
}

.error {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	color: red;
	margin-top: 0.5rem;
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
