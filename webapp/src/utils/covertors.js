// Function to parse a duration string (e.g., "1h 23m") to an integer in minutes
export function parseDurationToInt(durationStr) {
    const hoursMatch = durationStr.match(/(\d+)h/); // Matches hours part
    const minutesMatch = durationStr.match(/(\d+)m/); // Matches minutes part

    const hours = hoursMatch ? parseInt(hoursMatch[1], 10) : 0;
    const minutes = minutesMatch ? parseInt(minutesMatch[1], 10) : 0;

    return hours * 60 + minutes;
}

// Function to format an integer (in minutes) to a duration string (e.g., "1h 23m")
export function formatDuration(totalMinutes) {
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;

    return hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
}


export function formatToFancyDateString(datetimeString) {
    const date = new Date(datetimeString);

    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    return date.toLocaleDateString(undefined, options);
}
