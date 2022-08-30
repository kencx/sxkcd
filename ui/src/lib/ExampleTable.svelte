<script>
	let examples = [
		{ query: "foo|bar", description: "foo or bar", },
		{ query: "-foo", description: "Exclude foo", },
		{ query: "foo*", description: "Any words that begin with foo", },
		{ query: "#420", description: "Filter by comic number", },
		{ query: "#420-690", description: "Between number range", },
		{ query: "@date: 2022-01-19", description: "From date to present", },
		{ query: "@date: 2022-01-01, 2022-08-01", description: "Between date range", },
	]
</script>

<div class="examples">
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

<style>
.examples details {
	font-size: 0.8rem;
	border-bottom: none;
	width: 80%;
	margin: auto;
}

.examples table {
	margin-bottom: 0%;
}

.examples th {
	font-size: 0.8rem;
}

@media (min-width:320px) and (max-width:640px) {
	.examples details {
		width: 100%;
		margin: none;
	}
}
</style>
