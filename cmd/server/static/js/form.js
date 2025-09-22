document.getElementById("pageTitle").innerHTML = "Update Your Details | Fidelis College"

const primaryForm = document.getElementById('primaryContactForm');
const secondaryForm = document.getElementById('secondaryContactForm');

const forms = [primaryForm,secondaryForm]

function showContactCards() {
    let contactCards = document.getElementById("contact-cards")
    let backBtn = document.getElementById("backBtn")
    let successMsg = document.getElementById("success-message")
    let activeForm;

    for (const form of forms) {
        if (form.style.display === "") {
            activeForm = form
        }
    }

    if (activeForm) {
        activeForm.style.display = "none"
    }

    backBtn.style.display = "none"
    successMsg.style.display = "none"

    contactCards.style.display = ""
}

function showContactForm(formName) {
    // Your logic to show the appropriate form
    // This could be:
    // 1. Redirect to form page with contact ID: window.location.href = `/form?contact=${contactNumber}`
    // 2. Show/hide form sections on the same page
    // 3. AJAX load the form data

    let chosenForm = document.getElementById(formName)
    let contactCards = document.getElementById("contact-cards")
    let backBtn = document.getElementById("backBtn")

    contactCards.style.display = "none"

    chosenForm.style.display = ""
    backBtn.style.display = "inline-flex"
}

document.addEventListener('DOMContentLoaded', function() {
    const successMessage = document.getElementById('success-message');
    const selectContactSubtitle = document.getElementById("selectContactSubtitle")

    const primaryCard = document.getElementById("primaryCard")
    const primaryCardName = document.getElementById("primaryCardInitial").innerHTML

    const secondaryCard = document.getElementById("secondaryCard")
    const secondaryCardName = document.getElementById("secondaryCardInitial").innerHTML

    if(primaryCardName === "") {
        primaryCard.style.display = "none"
    }

    if(secondaryCardName === "") {
        secondaryCard.style.display = "none"
    }

    if (primaryCardName === "" && secondaryCardName === "") {
        selectContactSubtitle.innerHTML="No contact information found, please contact data@fidelis.org.uk"
    }

    for (const form of forms) {
        form.style.display = 'none'

        form.addEventListener('submit', async function(event) {
            // Prevent the default form submission (which causes a redirect)
            event.preventDefault();

            // Get the form data
            const formData = new FormData(form);

            formData.set("contactType", form.id)

            // (Optional) Show a loading state on the button
            const submitButton = form.querySelector('button[type="submit"]');
            const originalButtonText = submitButton.textContent;
            submitButton.textContent = 'Saving...';
            submitButton.disabled = true;

            try {
                // Send the form data to your server endpoint
                const response = await fetch('/submit', {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    // Show the success message
                    successMessage.style.display = 'block';

                    // Optional: Hide the form after success
                    form.style.display = 'none';

                    // Optional: Scroll to the success message
                    successMessage.scrollIntoView({ behavior: 'smooth' });
                } else {
                    alert('There was a problem saving your details. Please refresh the page and try again.');
                }

            } catch (error) {
                console.error('Error:', error);
                alert('A network error occurred. Please refresh the page and try again.');
            } finally {
                // Re-enable the button regardless of success/failure
                submitButton.textContent = originalButtonText;
                submitButton.disabled = false;
            }
        });
    }
});

async function initMap(elementName) {
    // Request needed libraries.
    await google.maps.importLibrary("places");

    // Create the input HTML element, and append it.
    //@ts-ignore
    const placeAutocomplete = new google.maps.places.PlaceAutocompleteElement({
        includedRegionCodes: ['gb'],
    });

    placeAutocomplete.style.colorScheme = "light"
    placeAutocomplete.style.border = "1px solid #ced4da"
    placeAutocomplete.style.borderRadius = "0.375rem"
    //@ts-ignore
    document.getElementById(elementName).appendChild(placeAutocomplete);
    const subtitle = document.createElement("small");
    subtitle.style = "color: #6c757d; display: block; margin-top: 0.25rem;"
    subtitle.innerHTML = "Use the address search to update the fields below"
    document.getElementById(elementName).appendChild(subtitle)
    // Inject HTML UI.
    const selectedPlaceTitle = document.createElement('p');
    selectedPlaceTitle.textContent = '';
    document.body.appendChild(selectedPlaceTitle);
    const selectedPlaceInfo = document.createElement('pre');
    selectedPlaceInfo.textContent = '';
    document.body.appendChild(selectedPlaceInfo);
    // Add the gmp-placeselect listener, and display the results.
    //@ts-ignore
    placeAutocomplete.addEventListener('gmp-select', async ({ placePrediction }) => {
        const place = placePrediction.toPlace();
        await place.fetchFields({ fields: ['displayName', 'formattedAddress', 'location', 'postalAddress'] });
        console.log(JSON.stringify(place.toJSON(), /* replacer */ null, /* space */ 2))

        if (elementName==="addressSearchTwo") {
            for (let i = 0; i < 3; i++){
                document.getElementById("secondaryAddressLine"+i).value = ""
            }

            for (let i = 0; i < place.postalAddress.addressLines.length; i++) {
                document.getElementById("secondaryAddressLine"+i).value = place.postalAddress.addressLines[i]
            }
            document.getElementById("secondaryPostalCode").value = place.postalAddress.postalCode

            document.getElementById('secondaryCounty').value = place.postalAddress.locality;
        }

        for (let i = 0; i < 3; i++){
            document.getElementById("addressLine"+i).value = ""
        }

        for (let i = 0; i < place.postalAddress.addressLines.length; i++) {
            document.getElementById("addressLine"+i).value = place.postalAddress.addressLines[i]
        }
        document.getElementById("postalCode").value = place.postalAddress.postalCode

        document.getElementById('county').value = place.postalAddress.locality;

    });
}
initMap("addressSearch");
initMap("addressSearchTwo")