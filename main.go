package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION
// ============================================================================

const (
	GraphQLEndpoint = "https://bff-page.kakao.com/graphql"

	// Thay Cookie mới của bạn vào đây
	CookieString = "_kpdid=ca05be6623084b93b8b01c9e1200e0c3; _kpiid=8bc3c90e0a818233479d4791be94e941; _kau=f6d3196ee10463825a16bb96a31e6cd4feeac287925281bda834e2c715723b14f7ecbed5e08bdb56a3865e49effdd185ea76b19e08fbd409eca23545cf72065017bb053112c83a8fbedcca17caf64655d3f4792a53c4a5a9498d2c2e24d276a0b94871e2c9277fc4899ea615d57984a65d9709270690a2cc0120201e3634383033393236373737363634343533393233363338373838343837313733e5b0269151bb36ab059d321e6a892eb8; _kawlt=fdpHGu1r8jIb2SxETwVqbcl4_mZNnmLp7uXmlyS1IjR5m03-wJr77cWrS93KZg3pgyC09t-PvvMDJdBj8iAv7pWqsR4cYqe6dZzMsV8ucdqhrADx3RBLPdxrSlV_YyYQ; _kawltea=1765882153; _karmt=2b9aLc7eRdTsx-pWuZYERQWXpHk7ebyJqjI4ojnGcFQRDxPfO3Pi6cWfq7ySHnr8; _karmtea=1768409353; _kahai=ba4be1e91c915262b60a6402f5e1e321446b1a1739f6a79ad21e84f1393c8927; _kaslt=TXmgRsAl2CyiHK0xQwRK9+IhZtBaLNQDGdbzDahGjvEoJlHoZZaM1F+kMAoK/flKUUJjw3a10N8TYg6jYci9QA==; _kpawbat_e=YbFMiRZUgnFJodO2oY6lVPc04F2K%2F4hxEedpXyrKLoJgIyTLxwy%2F%2FAxNjkgaYGU%2Bc5VIUtvefZPnCsoy4XrkaIQ19VpK39jFwYBvpxXqQIQ%3D; _kpwtkn=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..qLhvAN2c4dxG8tSR.HrI6sY70hq7VYhskDwiNzb8hHsp1JrUf2QxuDUPTAp9_aSDVBOi-O_D7tokyWnswNLb_xcYmk4nuhhQvrwsGnz_7xH6Hz9RvJ2c6yHckX0-excd9Rp1RaF54wBECTN4IKN7A2zD4Hf4aKCRl69lQMk2B43T3dYfekoZVl58EEU9CEVclcbdK0nhuNCsVExEIpkTW1lnDmM1XA51D3t9Csjd3KOtSj9qBzxV9aZ4HPFtiBKoDYw0mXIf5wNdC6UqOjB954376Sb9No_HwW8bPQv_gTpyizwQhkjMhU6tKaqo7N3LmCTU7_ieCgo_xm3UhO9zAnQYh0A6zaXCbucwBSbAUMsNsLnpgxJvqlDlSjlFb-xfjUbRlXAXhDgmq9nFCaa6O00E2wUlxkEB16cYVXfI7XUIwCZ3cI3HNxFORdOHm_Hd6pSSFKQiaknQ3lfGvr7Cia89GLJGRNtVFHfd71N5lPYJUtiVbVrlI4gAIRGqpkET5LNLwE8HslpItm0FsPQBNnfBUOnwmUyuLVGS2OPUirLWzKcaTz5oWklxquS0mke-GdehBVWErjKCfKiEQolgqhP8DjnWsV6OJu1W-kahk1izJ-u8pcOpXt6YSlzr04rQp2N9Fr3GjQk5dFdp8BNiWhVHiU1euUkpc5hiqep7HYHkjqKnAhW7M6AXSCX-OpZF81x2OR71NAvkioJCpeXs3Y0goxiNCQTf5NX8UtjWHj8HxN5zgp4Q9dzQ6npNmw6tSE6sHzPCnBZWHdpTENAgUApE0A6c4gxbOLF2ryneKdFRoNlTbAZKF8vxs5Bdc8Oq8THSRFMZMQ6M9XMB7zPiLmbJwKukWjKFIkeU6djeTxCJIu-PE2BP7_m9fjWdcPkMKWIMuVjuc3U7BI3QmWIcfUPWuW98U4Z-pmZdg_DJMEt8pynAs1sFUvVNeEkKaqQO8-9wT4Kx3YJJ3fg2kOXq5pfu3iFSsPt4N7K0MC9YMUOFzvrAIO2WKgZdHWMl1OusfaKSk11ZNo3K1xEOWxMFytaOhWkE5yiUhKOlylNVKGie_RBcvXzGRD2M0pn9mv-z0VN1VVDFhhO45X3upaD4lohISI7DpN3lgwV8rZ0k_NRMZ6j29f7fxpjI23yvp5AmXmm3cgshKKjmW81KZ0bSpyhSxCfHjLN5Cf2hkpafR6QMHGcpBROZUXyVWK4SfcclkqODBhhVmU2kdHO8rNNB54ahachNAaghp9BfY-wyKXPZBu5NfLGhM7xwHZNU8OYs_zysr9l2_0v45L0Fksl8tTWw1K5FnfIH3K00jQdoNYU9s0VcZshXFPjKETXvDuFTEwmYhC86o90XMun2sT2cdmM931ygRGEfYdJbv-yZMF8smC2jPoEuszLozpKhPNzFlhsHJmjZZlnwSD2yYMhRpIOtWWheUfdIK0lRCH8eQDu3se4AUX2brxF4ct4yplxu6SMsRInKogVTC71hsudmMQmS3kmybjtuOuGCJKBOUPF9v7yRTjKOWMFoWVMCmtnPj_8WJ3Bw5scnN_BlErba9Q-9_oYQtTKde5Ot5WwQ19IF-e7vVSxqq0NJicOq5pChK7q0hTB_bzCo_Uu-rFVQIsB2p15HdhzDNl0PZQsrmzak2aERx2o06UxZ_ohBVUIFWdo9uhBGyubE77OM3i_txxN05-c38DCJ1ooFNRy0WTw6YpIYFtqXp1g-rNiAnwSuW_2KLzFj_taWmB4_NH3JedqO302USftRSdTAKtF_kmRzLfc3huI90Gib5HVPUOOhPY1u5Sc8QQsjQJELAbKh2vYgGTtb-6C2t3ex2kQve3fynjWH8S6bXTdthJFRVL20E92dcdCWbAj0ztvp4BHNpLVer6BkEOlymyVRse11M3o5K5GlqRrN-Gk7WiOzPpB8dfTqQdAPT-qMcBgC9ioZv7XG5FHDqDBxecW_u12ObcROr5RV6apH2hwJkqgDQrnCUXX8XZlfapU84BKioJDLlzEFJIqS4DKWH51qwDy-2PxCsl9pzUejIZfFli_LjZSM9W8Ut15u4-K-GQXB_2fcDkzeD-w5Q3kKCaFlHfUOYu0XAaSkIR_k8Xju7vDL3OUCDVA07aTRSgvcid4nh8nyVWZskAy-i1SEk4AE5MwpsHq1mHmuGXgI8kWrUVtTEttILDY8uelkD7lGzQeHhppOeJrWhWmTwGZzUdE4EtjJZqCsEHhfF_mi_LviWna1XcByd3r6Cv7SsPz4Qez0o6njkUYLPq5Vlc6WAL8Gv6oYOEohh-7JgeAdPppgu1QzdaRibcyuwFbNUGZTvGjOpncxO.S1qWC0AedTW22JIh9n2bnw; _T_ANO=KrZJsDoqUmhZvcqmtZtsu3wZm0lkkvcy5awtanxkn6B3AoTOC4fTdwn39W6FBnXDDw4PXjvAPrg2ueTf3yMRTOvujU8knsp513Qym4K/ckfsqc6+HT09iggvW3FJeCd58n0d8PkLsvwwwS/aZz+t3cFv1KZL4jRxlBJNIZaC+XmRfLFDH76Y0Nd1MZQXK0zUK3K0jTHWvJsE8kzuPnBwOyZpnjvnDLgqHdqlzBNo4HdKiWFIcNcTivzGwsXJ8GKF5S1dab2ue31WgwhKwciz+1W0ZVuRGij0Y2UZz3nhkKhSDLa195GCmA/xNVAIi74wZ/WLEPBUcTlJXzTLHrHjuQ=="
)

// ============================================================================
// QUERY (dùng payload bạn cung cấp — giữ nguyên các type Long!)
// ============================================================================

const ViewerQuery = `
query viewerInfo($seriesId: Long!, $productId: Long!) {
  viewerInfo(seriesId: $seriesId, productId: $productId) {
    item {
      ...SingleFragment
    }
    seriesItem {
      ...SeriesFragment
    }
    prevItem {
      ...NearItemFragment
    }
    nextItem {
      ...NearItemFragment
    }
    viewerData {
      ...TextViewerData
      ...TalkViewerData
      ...ImageViewerData
      ...VodViewerData
    }
    displayAd {
      ...DisplayAd
    }
  }
}

fragment SingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  series {
    ...SeriesFragment
  }
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  discountRate
  discountRateText
  isShortsDrama
}

fragment SeriesFragment on Series {
  id
  seriesId
  title
  thumbnail
  landThumbnail
  categoryUid
  lang
  category
  categoryType
  subcategoryUid
  subcategory
  badge
  isAllFree
  isWaitfree
  ageGrade
  state
  onIssue
  authors
  description
  pubPeriod
  freeSlideCount
  lastSlideAddedDate
  waitfreeBlockCount
  waitfreePeriodByMinute
  bm
  saleState
  startSaleDt
  saleMethod
  discountRate
  discountRateText
  serviceProperty {
    ...ServicePropertyFragment
  }
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
  translateProperty {
    ...TranslatePropertyFragment
  }
}

fragment ServicePropertyFragment on ServiceProperty {
  viewCount
  readCount
  ratingCount
  ratingSum
  commentCount
  pageContinue {
    ...ContinueInfoFragment
  }
  todayGift {
    ...TodayGift
  }
  preview {
    ...PreviewFragment
  }
  waitfreeTicket {
    ...WaitfreeTicketFragment
  }
  isAlarmOn
  isLikeOn
  ticketCount
  purchasedDate
  lastViewInfo {
    ...LastViewInfoFragment
  }
  purchaseInfo {
    ...PurchaseInfoFragment
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment ContinueInfoFragment on ContinueInfo {
  title
  isFree
  productId
  lastReadProductId
  scheme
  continueProductType
  hasNewSingle
  hasUnreadSingle
}

fragment TodayGift on TodayGift {
  id
  uid
  ticketType
  ticketKind
  ticketCount
  ticketExpireAt
  ticketExpiredText
  isReceived
  seriesId
}

fragment PreviewFragment on Preview {
  item {
    ...PreviewSingleFragment
  }
  nextItem {
    ...PreviewSingleFragment
  }
  usingScroll
}

fragment PreviewSingleFragment on Single {
  id
  productId
  seriesId
  title
  thumbnail
  badge
  isFree
  ageGrade
  state
  slideType
  lastReleasedDate
  size
  pageCount
  isHidden
  remainText
  isWaitfreeBlocked
  saleState
  operatorProperty {
    ...OperatorPropertyFragment
  }
  assetProperty {
    ...AssetPropertyFragment
  }
}

fragment OperatorPropertyFragment on OperatorProperty {
  thumbnail
  copy
  helixImpId
  isTextViewer
  selfCensorship
  isBook
  cashInfo {
    discountRate
    setDiscountRate
  }
  ticketInfo {
    price
    discountPrice
    ticketType
  }
}

fragment AssetPropertyFragment on AssetProperty {
  bannerImage
  cardImage
  cardTextImage
  cleanImage
  ipxVideo
  bannerSet {
    ...BannerSetFragment
  }
  cardSet {
    ...CardSetFragment
  }
  cardCover {
    ...CardCoverFragment
  }
}

fragment BannerSetFragment on BannerSet {
  backgroundImage
  backgroundColor
  mainImage
  titleImage
}

fragment CardSetFragment on CardSet {
  backgroundColor
  backgroundImage
}

fragment CardCoverFragment on CardCover {
  coverImg
  coverRestricted
}

fragment WaitfreeTicketFragment on WaitfreeTicket {
  chargedPeriod
  chargedCount
  chargedAt
}

fragment LastViewInfoFragment on LastViewInfo {
  isDone
  lastViewDate
  rate
  spineIndex
}

fragment PurchaseInfoFragment on PurchaseInfo {
  purchaseType
  rentExpireDate
  expired
}

fragment TranslatePropertyFragment on TranslateProperty {
  category {
    ...LocaleMapFragment
  }
  sub_category {
    ...LocaleMapFragment
  }
}

fragment LocaleMapFragment on LocaleMap {
  ko
  en
  th
}

fragment NearItemFragment on NearItem {
  productId
  slideType
  ageGrade
  isFree
  title
  thumbnail
}

fragment TextViewerData on TextViewerData {
  type
  atsServerUrl
  metaSecureUrl
  contentsList {
    chapterId
    contentId
    secureUrl
  }
}

fragment TalkViewerData on TalkViewerData {
  type
  talkDownloadData {
    dec
    host
    path
    talkViewerType
  }
}

fragment ImageViewerData on ImageViewerData {
  type
  imageDownloadData {
    ...ImageDownloadData
  }
}

fragment ImageDownloadData on ImageDownloadData {
  files {
    ...ImageDownloadFile
  }
  totalCount
  totalSize
  viewDirection
  gapBetweenImages
  readType
}

fragment ImageDownloadFile on ImageDownloadFile {
  no
  size
  secureUrl
  width
  height
}

fragment VodViewerData on VodViewerData {
  type
  vodDownloadData {
    contentPackId
    drmType
    endpointUrl
    width
    height
    duration
  }
  drmInfo {
    type
    serverType
    error
    fairplayCertificateUrl
    widevineLicenseUrl
    fairplayLicenseUrl
    token
    provider
    assertion
  }
}

fragment DisplayAd on DisplayAd {
  sectionUid
  bannerUid
  treviUid
  momentUid
}
`

// ============================================================================
// MODELS (mình chỉ khai báo những field cần dùng)
// ============================================================================

type GraphQLPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type GQLResponse struct {
	Data struct {
		ViewerInfo ViewerInfo `json:"viewerInfo"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ViewerInfo struct {
	Item       Item       `json:"item"`
	SeriesItem *Series    `json:"seriesItem,omitempty"`
	PrevItem   *NearItem  `json:"prevItem,omitempty"`
	NextItem   *NearItem  `json:"nextItem,omitempty"`
	ViewerData ViewerData `json:"viewerData"`
}

type Item struct {
	Title string `json:"title"`
	// thêm field nếu cần về sau
	ProductId int `json:"productId,omitempty"`
	SeriesId  int `json:"seriesId,omitempty"`
}

type Series struct {
	SeriesId int    `json:"seriesId"`
	Title    string `json:"title"`
}

type NearItem struct {
	ProductId int    `json:"productId"`
	Title     string `json:"title"`
	// slideType, etc. omitted — thêm nếu cần
}

type ViewerData struct {
	AtsServerUrl string        `json:"atsServerUrl"`
	MetaSecure   string        `json:"metaSecureUrl,omitempty"`
	ContentsList []ContentPart `json:"contentsList"`
}

type ContentPart struct {
	ChapterId int    `json:"chapterId,omitempty"`
	ContentId int    `json:"contentId,omitempty"`
	SecureUrl string `json:"secureUrl"`
}

type ContentFileResponse struct {
	ContentInfo struct {
		ParagraphList []ParagraphNode `json:"paragraphList"`
	} `json:"contentInfo"`
}

type ParagraphNode struct {
	Type     string          `json:"type"`
	Text     string          `json:"text"`
	Children []ParagraphNode `json:"childParagraphList,omitempty"`
}

// ============================================================================
// MAIN: loop theo nextItem.productId — giữ flow tải/chạy y nguyên code gốc
// ============================================================================

func main() {
	// context tổng cho toàn bộ chạy
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	headers := map[string]string{
		"accept":          "*/*",
		"accept-language": "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7",
		"content-type":    "application/json",
		"origin":          "https://page.kakao.com",
		"referer":         "https://page.kakao.com/",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"cookie":          CookieString,
	}

	// Bắt đầu từ đây: set seriesId và productId chương ban đầu bạn muốn crawl
	currentSeries := 61097897
	currentProduct := 65061307

	for {
		fmt.Printf("\n>>> Crawl productId = %d (seriesId=%d)\n", currentProduct, currentSeries)

		// 1) build payload giữ nguyên như bạn yêu cầu
		payload := GraphQLPayload{
			OperationName: "viewerInfo",
			Query:         ViewerQuery,
			Variables: map[string]interface{}{
				"seriesId":  currentSeries,
				"productId": currentProduct,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		respBody, err := makeRequest(ctx, client, "POST", GraphQLEndpoint, headers, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf(">>> ERROR: GraphQL request failed: %v\n", err)
			break
		}

		var gqlResp GQLResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			fmt.Printf(">>> ERROR: Parse Error: %v\nBody: %s\n", err, string(respBody))
			break
		}
		if len(gqlResp.Errors) > 0 {
			fmt.Printf(">>> ERROR: GraphQL Error: %s\n", gqlResp.Errors[0].Message)
			break
		}

		viewerData := gqlResp.Data.ViewerInfo.ViewerData
		itemInfo := gqlResp.Data.ViewerInfo.Item
		nextItem := gqlResp.Data.ViewerInfo.NextItem

		if viewerData.AtsServerUrl == "" {
			fmt.Println(">>> ERROR: Không tìm thấy atsServerUrl. Kiểm tra lại Cookie.")
			break
		}

		// Tạo file
		safeTitle := sanitizeFilename(itemInfo.Title)
		fileName := safeTitle + ".txt"

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf(">>> ERROR: Không thể tạo file %s: %v\n", fileName, err)
			break
		}
		fmt.Printf(">>> Đã tạo file: %s\n", fileName)
		fmt.Printf(">>> Tổng số phần: %d. Bắt đầu tải...\n", len(viewerData.ContentsList))

		// Ghi tiêu đề
		if _, err := file.WriteString(fmt.Sprintf("# %s\n\n", itemInfo.Title)); err != nil {
			fmt.Printf(">>> WARNING: lỗi ghi tiêu đề: %v\n", err)
		}

		// Loop tải các phần
		for i, part := range viewerData.ContentsList {
			if i > 0 {
				time.Sleep(1 * time.Second)
			}
			fullUrl := viewerData.AtsServerUrl + part.SecureUrl
			fmt.Printf(" -> [%d/%d] Đang tải... (Delay 1s)\n", i+1, len(viewerData.ContentsList))

			partBody, err := makeRequest(ctx, client, "GET", fullUrl, headers, nil)
			if err != nil {
				fmt.Printf("   [Lỗi tải phần %d]: %v\n", i+1, err)
				continue
			}

			var contentResp ContentFileResponse
			if err := json.Unmarshal(partBody, &contentResp); err != nil {
				// nếu không parse được JSON -> có thể là image / binary -> bỏ qua
				continue
			}

			for _, node := range contentResp.ContentInfo.ParagraphList {
				text := extractTextRecursive(node)
				text = strings.ReplaceAll(text, "&nbsp;", " ")
				if _, err := file.WriteString(text); err != nil {
					fmt.Printf("   Lỗi ghi file: %v\n", err)
				}
			}
		}

		file.Close()
		fmt.Println(">>> XONG! Kiểm tra file:", fileName)

		// Nếu không có nextItem -> dừng
		if nextItem == nil || nextItem.ProductId == 0 {
			fmt.Println(">>> HẾT CHƯƠNG hoặc không tìm thấy nextItem.")
			break
		}

		// Cập nhật productId để crawl chương tiếp theo
		fmt.Printf(">>> Next: %d (%s)\n", nextItem.ProductId, nextItem.Title)
		currentProduct = nextItem.ProductId

		// delay an toàn trước lần tiếp theo
		time.Sleep(2 * time.Second)
	}
}

// ============================================================================
// HELPERS
// ============================================================================

func makeRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %d | Body: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}

func extractTextRecursive(node ParagraphNode) string {
	var sb strings.Builder
	if node.Text != "" {
		sb.WriteString(node.Text)
	}
	if node.Type == "BR" {
		sb.WriteString("\n")
	}
	for _, child := range node.Children {
		sb.WriteString(extractTextRecursive(child))
	}
	switch node.Type {
	case "P", "DIV":
		sb.WriteString("\n")
	case "H3":
		sb.WriteString("\n\n### ")
	}
	return sb.String()
}

func sanitizeFilename(name string) string {
	reg, _ := regexp.Compile(`[\\/:*?"<>|]`)
	safe := reg.ReplaceAllString(name, "_")
	return strings.TrimSpace(safe)
}
