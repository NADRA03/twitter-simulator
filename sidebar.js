export function render() {
    return `
        <img src="/assets/logo.svg" class="logo" />
        <nav class="nav1">
            <ul>
<li>
<a href="/" onclick="loadPage('home'); return false;">
Home
</a>
</li>
                <li><a href="/log-in" onclick="loadPage('log-in'); return false;">Log-in</a></li>
                <li><a href="/chats" onclick="loadPage('chats'); return false;">Chats</a></li>
                <li><a href="/search" onclick="loadPage('search'); return false;">Search</a></li>
                <li><a href="/profile" onclick="loadPage('profile'); return false;">Profile</a></li>
            </ul>
                                <form id="logoutForm" action="/logout" method="POST">
        <button id="logoutButton" type="submit">Logout</button>
    </form>
        </nav>
    `;
}