document.getElementById("pageTitle").innerHTML = "Update Your Details | Fidelis College"

const forms = document.getElementsByTagName("form")

let personIndex = -1;
let children;
let editingIndex = -1;

function showContactCards() {
    let contactCards = document.getElementById("contact-cards")
    let backBtn = document.getElementById("backBtn")
    let successMsg = document.getElementById("success-message")
    let selectLang = document.getElementById("ealSpoken").children
    let activeForm;

    for (const form of forms) {
        if (form.style.display === "" && form.id !== "conditionForm") {
            activeForm = form
        }
    }

    if(!activeForm || confirm("Are you sure you want to go back? All unsaved changes will be lost.")) {
        for (const option of selectLang) {
            option.selected = false
        }

        if (activeForm) {
            activeForm.style.display = "none"
        }

        backBtn.style.display = "none"
        successMsg.style.display = "none"

        contactCards.style.display = ""
    }
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

function showChildForm(schoolId) {
    let childForm = document.getElementById("childForm")

    let medConsentSuffixId = "-med-consent"
    let photoConsentSuffixId ="-photo-consent"

    personIndex = schoolId

    fetch('/children').then(async (resp) => {
        children = await resp.json()

        if(children[personIndex]["photoConsent"]) {
            document.getElementById("yes-photo-consent").checked = true
        } else {
            document.getElementById("no-photo-consent").checked = true
        }

        if(children[personIndex]["medConsent"]) {
            document.getElementById("yes-med-consent").checked = true
        } else {
            document.getElementById("no-med-consent").checked = true
        }

        if(children[personIndex]["tripsConsent"]) {
            document.getElementById("yes-trips-consent").checked = true
        } else {
            document.getElementById("no-trips-consent").checked = true
        }

        if(children[personIndex]["isEal"]) {
            document.getElementById("yes-eal").checked = true
        } else {
            document.getElementById("no-eal").checked = true
        }

        children[personIndex]["toDelete"] = []

        for (const lang of children[personIndex]["languages"]) {
            document.getElementById(lang).selected = true
        }

        renderConditionsTable()
        showContactForm("childForm")
    }).catch(err => {
        console.log(err)
    })
}

document.addEventListener('DOMContentLoaded', async function() {
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
        if(form.id !== "conditionForm") {
            form.style.display = "none"
        }

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
                let response;
                if(form.id === "childForm") {
                    let formPhotoConsent = formData.get("photoConsent") === "yes"
                    let formMedConsent = formData.get("medConsent") === "yes"
                    let isEal = formData.get("isEal") === "yes"
                    let tripsConsent = formData.get("tripsConsent") === "yes"

                    let languagesSelected = formData.getAll("ealSpoken")

                    children[personIndex]["medConsent"] = formMedConsent;
                    children[personIndex]["photoConsent"] = formPhotoConsent;
                    children[personIndex]["tripsConsent"] = tripsConsent;
                    children[personIndex]["languages"] = languagesSelected;
                    children[personIndex]["isEal"] = isEal;
                    children[personIndex]["schoolId"] = personIndex.toString()


                    response = await fetch("/updateChildren", {
                        method: 'POST',
                        body: JSON.stringify(children[personIndex])
                    })
                } else {
                    response = await fetch("/submit", {
                        method: 'POST',
                        body: formData
                    });
                }

                if (response.ok) {
                    // Show the success message
                    successMessage.style.display = 'block';

                    // Optional: Hide the form after success
                    form.style.display = 'none';

                    // Optional: Scroll to the success message
                    successMessage.scrollIntoView({ behavior: 'smooth' });
                } else {
                    alert("Error: "+await response.text());
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


function openConditionModal(editIndex = -1) {
    editingIndex = editIndex;
    const modal = document.getElementById('conditionModal');
    const form = document.getElementById('conditionForm');

    if (editIndex >= 0) {
        // Editing existing condition
        const condition = children[personIndex]["conditions"][editIndex];
        document.getElementById('conditionSelect').value = getConditionType(condition.groupId);
        document.getElementById('medicationName').value = condition.treatment;
        document.getElementById('conditionName').value = condition.type;
        document.getElementById('conditionNotes').value = condition.furtherInfo || '';
        document.querySelector('.modal-header h3').textContent = 'Edit Condition';
        document.querySelector('.btn--primary').textContent = 'Update Condition';
    } else {
        // Adding new condition
        form.reset();
        document.querySelector('.modal-header h3').textContent = 'Add Condition';
        document.querySelector('.btn--primary').textContent = 'Add Condition';
    }

    modal.style.display = 'block';
}

function closeConditionModal() {
    document.getElementById('conditionModal').style.display = 'none';
    editingIndex = -1;
}

function addCondition() {
    const conditionTypeId = document.getElementById('conditionSelect').selectedOptions;
    const conditionName = document.getElementById("conditionName").value;
    const medication = document.getElementById('medicationName').value;
    const notes = document.getElementById('conditionNotes').value;

    if (!conditionTypeId || !medication || !conditionName) {
        alert('Please provide a condition type, name and medication');
        return;
    }

    const conditionData = {
        groupId: parseInt(conditionTypeId[0].id.substring(8)),
        type: conditionName,
        treatment: medication,
        furtherInfo: notes,
        toBeUploaded: true,
    };

    if (editingIndex >= 0) {
        // Update existing condition
        let condition = children[personIndex]["conditions"][editingIndex]

        condition.groupId = conditionData.groupId;
        condition.type = conditionData.type;
        condition.treatment = conditionData.treatment;
        condition.furtherInfo = conditionData.furtherInfo;
        condition.toBeUploaded = true;
    } else {
        // Add new condition
        children[personIndex]["conditions"].push(conditionData);
    }

    renderConditionsTable();
    closeConditionModal();

    // Update hidden form field for submission
}

function deleteCondition(index) {
    if (confirm('Are you sure you want to delete this condition?')) {
        let toDelete = children[personIndex]["conditions"].splice(index, 1)[0];
        children[personIndex]["toDelete"].push(toDelete)
        renderConditionsTable();
    }
}

function renderConditionsTable() {
    const tbody = document.getElementById('conditionsBody');
    const emptyState = document.getElementById('emptyState');
    const table = document.getElementById('table')

    tbody.innerHTML = '';

    if (children[personIndex]["conditions"].length === 0) {
        emptyState.classList.remove('hidden');
        table.classList.remove('show');
    } else {
        emptyState.classList.add('hidden');
        table.classList.add('show');

        children[personIndex]["conditions"].forEach((condition, index) => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${getConditionType(condition.groupId)}</td>
                <td>${condition.type}</td>
                <td>${getMedicationDisplayName(condition.treatment)}</td>
                <td>${condition.furtherInfo || '-'}</td>
                <td class="actions">
                    <button type="button" class="btn--delete" onclick="deleteCondition(${index})">Delete</button>
                    <button type="button" class="btn--edit" onclick="openConditionModal(${index})">Edit</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    }
}

function getConditionType(value) {
    let condType = document.getElementById("condType"+value)

    return condType.innerHTML
}

function getMedicationDisplayName(value) {
    const names = {
        'ventolin': 'Ventolin',
        'insulin': 'Insulin',
        'epipen': 'EpiPen',
        'ritalin': 'Ritalin',
        'other': 'Other'
    };
    return names[value] || value;
}


// Close modal when clicking outside
window.addEventListener('click', function(event) {
    const modal = document.getElementById('conditionModal');
    if (event.target === modal) {
        closeConditionModal();
    }
});