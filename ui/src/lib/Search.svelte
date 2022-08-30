<script lang="ts">
	import Comic from '$lib/Comic.svelte';
	import { search } from './search'

	let query: string = '';
	let page: number = 1;
	let timer: any;

	let examples = [
		{ query: "foo|bar", description: "foo or bar", },
		{ query: "-foo", description: "Exclude foo", },
		{ query: "foo*", description: "Any words that begins with foo", },
		{ query: "#420", description: "Filter by comic number", },
		{ query: "#420-690", description: "Between number range", },
		{ query: "@date: 2022-01-19", description: "From date to present", },
		{ query: "@date: 2022-01-01, 2022-08-01", description: "Between date range", },
	]

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
	<div class="syntax">
		<details>
			<summary>Examples</summary>
			<table>
				<thead>
					<tr>
						<th>Query</th>
						<th>Description</th>
					</tr>
				</thead>
				<tbody>
					{#each examples as example}
					<tr>
						<th><code>{example.query}</code></th>
						<th>{example.description}</th>
					</tr>
					{/each}
				</tbody>
			</table>
		</details>
	</div>

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
		<small class="search-message error">An error occured: {error.message}</small>
	{/await}
</div>

<style>
.syntax details {
	font-size: 0.8rem;
	border-bottom: none;
	width: 80%;
	margin: auto;
}

.syntax table {
	margin-bottom: 0%;
}

.syntax th {
	font-size: 0.8rem;
}

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
	color: red;
	margin-top: 0.5rem;
	margin: auto;
}

@media (min-width:320px) and (max-width:640px) {
	.syntax details {
		width: 100%;
		margin: none;
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
