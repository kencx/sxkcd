
import type Comic from '$lib/Comic.svelte';

export type SearchResult = {
	count: number;
	comics: Comic[];
	time: number;
};

function sanitize(input: string): string {
	//TODO
	return ""
}

function extractDateRange(input: string): string[] | null {
	var rx = /@date:\s?([0-9]{4}\-[0-9]{2}\-[0-9]{2})\s?[,]?\s?([0-9]{4}\-[0-9]{2}\-[0-9]{2})?/g;
	var regexp = new RegExp(rx);
	var matches = regexp.exec(input);
	return matches
}

function epoch(dateStr: string): number {
	var ms = new Date(dateStr).getTime();
	return Math.floor(ms / 1000);
}

export async function search(query: string, page: number): Promise<SearchResult | null> {
	if (query == "") {
		return null;
	}

	// query = sanitize(query)

	// handle date query
	if (query.includes("@date:")) {
		var match = extractDateRange(query);

		if (match != null) {
			var from = epoch(match[1]);
			var to = (match[2] == undefined) ? Math.floor(Date.now() / 1000) : epoch(match[2]);

			if (!isNaN(from) && !isNaN(to)) {
				query = `@date:[${from} ${to}]`;
			} else {
				throw Error("invalid date format");
			}
		} else {
			throw Error("invalid date format");
		}
	}

	let abort = new AbortController();
	try {
		var endpoint = `/search?q=${query}&page=${page}`;
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

