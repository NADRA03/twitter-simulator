export function render() {
    return `
<svg class="logo" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Layer_1" width="24px" height="24px" viewBox="0 0 24 24" style="enable-background:new 0 0 24 24;" xml:space="preserve"><path d="M14.095479,10.316482L22.286354,1h-1.940718l-7.115352,8.087682L7.551414,1H1l8.589488,12.231093L1,23h1.940717  l7.509372-8.542861L16.448587,23H23L14.095479,10.316482z M11.436522,13.338465l-0.871624-1.218704l-6.924311-9.68815h2.981339  l5.58978,7.82155l0.867949,1.218704l7.26506,10.166271h-2.981339L11.436522,13.338465z"/></svg>
        <nav class="nav1">
            <ul>
                <li><a href="/" onclick="loadPage('home'); return false;">Home</a></li>
                <li><a href="/sign-up" onclick="loadPage('sign-up'); return false;">Sign Up</a></li>
                <li><a href="/log-in" onclick="loadPage('log-in'); return false;">Log-in</a></li>
                <li><a href="/log-in" onclick="loadPage('ntify'); return false;">Notifications</a></li>
                <li><a href="/chats" onclick="loadPage('chats'); return false;">Chats</a></li>
                <li><a href="/profile" onclick="loadPage('profile'); return false;">Profile</a></li>
                <!-- Add more chat links as needed -->
            </ul>
        </nav>
    `;
}