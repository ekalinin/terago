// Modal functionality using D3.js
var modal = d3.select("#descriptionModal");
var modalTitle = d3.select("#modalTitle");
var modalDescription = d3.select("#modalDescription");
var closeBtn = d3.select(".close");

// Create a map for faster lookup by ID
var entriesById = {};
radarEntries.forEach(function(entry, index) {
    // The radar library assigns IDs starting from 1
    entriesById[index + 1] = entry;
});

// Function to show modal
function showModal(entry) {
    if (entry && entry.description) {
        modalTitle.text(entry.label);
        modalDescription.text(entry.description);
        modal.style("display", "block");
    }
}

// Function to hide modal
function hideModal() {
    modal.style("display", "none");
}

// Wait for radar to be rendered and force simulation to complete
setTimeout(function() {
    // Add click handlers to all blip groups
    d3.selectAll("#radar g.blip")
        .style("cursor", "pointer")
        .on("click", function(event, d) {
            event.preventDefault();
            event.stopPropagation();

            // The data is bound to the element, so we can access it directly
            if (d && d.label) {
                var entry = radarEntries.find(function(e) {
                    return e.label === d.label;
                });
                if (entry) {
                    showModal(entry);
                }
            }
        });

    // Add click handlers to legend items
    d3.selectAll("#radar .legend .blip-list-item")
        .style("cursor", "pointer")
        .on("click", function(event) {
            event.preventDefault();
            event.stopPropagation();

            var listItem = d3.select(this);
            var linkElement = listItem.select("a");

            if (linkElement.node()) {
                var fullText = linkElement.text().trim();
                // Extract label by removing the number prefix (e.g., "1. React" -> "React")
                var label = fullText.replace(/^\d+\.\s*/, "");
                
                var entry = radarEntries.find(function(e) {
                    return e.label === label;
                });

                if (entry) {
                    showModal(entry);
                }
            }
        });
}, 2000);

// Close modal when clicking the X button
closeBtn.on("click", function(event) {
    event.preventDefault();
    hideModal();
});

// Close modal when clicking outside of it
modal.on("click", function(event) {
    if (event.target === modal.node()) {
        hideModal();
    }
});

// Close modal on Escape key
d3.select("body").on("keydown", function(event) {
    if (event.key === "Escape" || event.keyCode === 27) {
        hideModal();
    }
});

