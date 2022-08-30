
import type Comic from '$lib/Comic.svelte';

export type SearchResult = {
	count: number;
	comics: Comic[];
	time: number;
};

export async function search(query: string): Promise<SearchResult | null> {
	if (query == "") {
		return null;
	}
	query = encodeURIComponent(query)

	let abort = new AbortController();
	try {
		var endpoint = `/search?q=${query}`;
		const response = await fetch(endpoint, {
			signal: abort.signal
		});

		if (!response.ok) {
			throw Error(`${response.status} ${response.statusText}`);
		}
		const result = await response.json();
		return {
			count: result.count,
			comics: result.results,
			time: result.query_time,
		};
	} catch (err: any) {
		if (err.name == 'AbortError') {
			console.log('fetch aborted');
			return null;
		} else if (err.message.includes("Syntax error")) {
			throw Error(`invalid query format`);
		} else {
			console.error(err);
			throw Error(`something went wrong`);
		}
	}
};

