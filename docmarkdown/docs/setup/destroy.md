# Destroying a Layer0 instance

During testing or migration, you may find that you need to delete a non-functional or outdated instance of Layer0. This section provides procedures for destroying (deleting) a Layer0 instance.

## Part 1: Clean Up Your Layer0 Environments
In order to destroy a Layer0 instance, you must first delete all environments in the instance.

**To delete Layer0 environments:**
<ol>
  <li>At the command prompt, type the following command to see a list of environments in your Layer0 instance:
    <ul>
      <li class="command">**l0 environment list**</li>
    </ul>

  <li>For each environment listed in the previous step, with the exception of the environments that begin with "api", issue the following command (where _environmentName_ is the name of the environment you want to delete):
    <ul>
      <li class="command">**l0 environment delete** _environmentName_ **--wait**</li>
    </ul>
  Repeat this step until all of the environments, with the exception of the "api" environments, have been deleted.</li>

  <li>At the command prompt, type the following command to list all of the certificates that exist in your Layer0 instance:
    <ul>
      <li class="command">**l0 certificate list**</li>
    </ul>
  </li>
  <li>For each certificate listed in the previous step, type the following command to delete the certificate (replacing _certificateName_ with the name of the certificate you want to delete):
    <ul>
      <li class="command">**l0 certificate delete** _certificateName_ **--wait**</li>
    </ul><br />
    Repeat this step until all of the certificates have been deleted. When you have finished deleting the environments and certificates in your Layer0 instance, proceed to Part 2.
  </li>
</ol>

## Part 2: Destroy the Layer0 instance

Once you have prepared your Layer0 instance for deletion, you can use the **l0-setup destroy** command to destroy the instance.

**To destroy a Layer0 instance:**

<ol>

  <li>At the command prompt, type the following command, replacing <em>prefixName</em> with the prefix you created when you created your Layer0 instance:
    <ul>
      <li class="command"><strong>l0-setup destroy</strong> <em>prefixName</em></li>
    </ul>
  </li>
</ol>
<div class="admonition note">
  <p class="admonition-title">Note</p>
  <p>The <strong>l0-setup destroy</strong> operation is idempotent (that is, it has no additional effects if you execute it multiple times with the same parameters). Therefore, if the <strong>destroy</strong> operation fails, you may be able to make it complete by running it again. If the <strong>destroy</strong> operation continues to fail after running it again, please contact the Xfra team at <strong>xfra@us.imshealth.com</strong>.</p>
</div>
