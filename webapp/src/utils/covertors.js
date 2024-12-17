// Function to parse a duration string (e.g., "1h 23m") to an integer in minutes
export function parseDurationToInt(durationStr) {
    const hoursMatch = durationStr.match(/(\d+)h/); // Matches hours part
    const minutesMatch = durationStr.match(/(\d+)m/); // Matches minutes part

    const hours = hoursMatch ? parseInt(hoursMatch[1], 10) : 0;
    const minutes = minutesMatch ? parseInt(minutesMatch[1], 10) : 0;

    return hours * 60 + minutes;
}

// Function to format an integer (in minutes) to a duration string (e.g., "1h 23m")
export function formatDuration(minutes) {
    const hours = Math.floor(minutes / 60);
    const remainingMinutes = minutes % 60;
    if (remainingMinutes === 0) {
        return `${hours}h`;
    }
    return `${hours > 0 ? `${hours}h ` : ''}${remainingMinutes}m`;
}


export function formatToFancyDateString(datetimeString) {
    const date = new Date(datetimeString);

    const options = { year: 'numeric', month: 'short', day: 'numeric' };
    return date.toLocaleDateString(undefined, options);
}

export const formatTime = (time) => {
    const date = new Date(time);

    let hours = date.getUTCHours();
    const minutes = date.getUTCMinutes();
    const ampm = hours >= 12 ? 'PM' : 'AM';

    // Convert to 12-hour format
    hours = hours % 12;
    hours = hours ? hours : 12; // The hour '0' should be '12'

    // Format minutes to always have two digits
    const minutesStr = minutes < 10 ? '0' + minutes : minutes;

    return `${hours}:${minutesStr} ${ampm}`;
};
