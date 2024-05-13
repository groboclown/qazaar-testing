# Use Case: Adding A New Feature

## Setup

Let's say we have a 3-tier web site (HTML + JavaScript user interface, a REST API server, and a database).  It has an existing feature for a User Profile.

* Feature: User Profile
    * Description: Users can edit their information, and other users may see a high level summary of other users' information.
    * Contains fields:
        - Name: User Name
          Type: User editable text field
          Editable by: Owning user
          Visible by: Any authenticated user
          Field size: 120 characters
        - Name: User Access Rights
          Type: access rights collection
          Editable by: administrator
          Visible by: Owning user, administrator

This gives us a high level overview of the expected behavior of the user profile.  It implies an ability to view and edit based on user permissions.

## Adding A Feature

Let's say that the development group has an update to the features come in to allow a user to edit their profile to include a bit of flavor text.  This allows users to get to know each other.  The user profile feature now includes:

* Feature: User Profile
    * Contains fields:
        - Name: Flavor Text

The requirements part of the project includes some Qazaar rules that force conformity from the requirements.  Processing the newly updated requirements through these rules causes Qazaar to flag the "Flavor Text" as a field but missing a field type label and authorization rules.  The product owner updates the profile to:

* Contains fields:
    - Name: Flavor Text
      Type: User editable text area
      Editable by: Owning user
      Visible by: Any authenticated user

Qazaar flags "Flavor Text" again, this time with missing the field size label, as all "text field" records need this.  So, the product owner makes another update:

* Contains fields:
    - Name: Flavor Text
      Type: User editable text area
      Editable by: Owning user
      Visible by: Any authenticated user
      Field size: 4,000 characters

## Trickle Down

With the new requirements in hand, the development team takes a pass at the project.  Qazaar makes some notes:

* The User Interface has a View User Profile (linked to the User Profile feature), but it does not include the "Flavor Text" field.
* The User Interface has an Edit User Profile, but it does not include the "Flavor Text" field.
* The REST API has "GET /user-profile/(name)" endpoint that returns the User Profile, but it does not include the world-readable "Flavor Text" field.
* The REST API has "PUT /user-profile/(name)" endpoint that allows updating the User Profile, but it does not include the user-editable "Flavor Text" field.
* The SQL Database Schema includes a User Profile table, but does not include the "Flavor Text" field.

The team also notes that, even with how amazing Qazaar is, it missed the REST API "PUT" support for individual fields within the user-profile, so they add in:

* The REST API does not include "PUT /user-profile/(name)/flavor-text".  Calls must have "Owning user" permissions.

The team reviews these requirements, and decides that this turns into changes to the existing product:

* User Interface Updates:
    * Add a new Flavor Text text area for Edit Profile actions.
    * Include the Flavor Text in the data sent to the server on a save action.
    * Add a new Flavor Text display area to the View Profile page.
* Server Updates:
    * Receive the Flavor Text field from Update Profile requests.
    * Send the Flavor Text field on Get Profile requests.
* Database Updates:
    * Add a new text blob column to the User Profile table.

Based on the feature and the requirements, the development team codifies these with:

* Field: User Profile - Flavor Text
    * Parent Record: UserProfile/username
    * Data type: text area
    * Editable: true
    * Maximum size: 4,000
    * Edit Access: owner
    * Read Access: any authorized user
    * Edit Label Name: "Describe Yourself"
    * View Label Name: "About Myself"

However, the product already has existing, *cross-cutting* requirements across the product:

* User Interface:
    * All text areas must include a label, whose value must pull from the supported localized translations.
        * Codified as: If has "Label", then Value must include localization.
    * All editable text fields must support Unicode characters.
        * Codified as: If "Data type" = "text area", then content = Unicode
    * All text areas must contain zero or more characters, and fewer than *per-field maximum* characters.
        * Codified as:
            * If "Data type" = "text area", then field must include 'maximum size' > 0.
            * If "Data type" = "text area", then UI field includes maximum length = field 'maximum size'.
    * The user interface must scrub the displayed text to not allow for cross-site scripting (XSS) vulnerabilities.
        * Codified as: If "Data type" = "text area", then UI escapes characters safe for rendering.
* Server:
    * All text fields received must use UTF-8 encoding.
    * The text fields must not exceed *per-field maximum* characters.
    * All text sent to the database must use proper parameter SQL commands to avoid SQL injection attacks.
    * Requests for saving data require *per-request level* authorized access.
    * Requests for sending data require *per-request level* authorized access.
    * Requests must respond within 2 seconds under heavy load.
* Database:
    * Columns storing more than 2000 Unicode characters must use binary blobs, encoded as UTF-32BE.
    * Binary blobs must set *per-field maximum* byte count, times 4 (because 4 bytes store 1 32-bit character).

On top of this, the development team has a series of items that it created because of discovered bugs and limitations in the software they depend on:

* All Unicode text must support 32-bit characters, such as emoji.
    * Codified as: If content = Unicode, then character range = valid unicode characters. (You can find the precise allowable values elsewhere; it's complex)
* All Unicode text can support up to 4 diacritic marks per displayed glyph.  The base character counts as one character, and each diacritic mark counts as another.  Beyond 4, the system strips them.

From these various data sources, we can now add meta-data to our tests to see what our tests cover.

