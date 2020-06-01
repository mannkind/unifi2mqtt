namespace Unifi.Models.Shared
{
    /// <summary>
    /// The shared key info => slug mapping across the application
    /// </summary>
    public class SlugMapping
    {
        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string MACAddress { get; set; } = string.Empty;

        /// <summary>
        /// 
        /// </summary>
        /// <value></value>
        public string Slug { get; set; } = string.Empty;

        /// <inheritdoc />
        public override string ToString() => $"Mac: {this.MACAddress}, Slug: {this.Slug}";
    }
}
